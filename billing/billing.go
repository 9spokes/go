package billing

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/9spokes/go/http"
	"github.com/9spokes/go/logging/v3"
	"github.com/go-redis/redis"
)

const StripeURL = "https://api.stripe.com"

type Context struct {
	APIKey      string
	CallbackURL string
	Cache       *redis.Client
}

type Profile struct {
	Name  string
	Email string
	ID    string
}

func (ctx *Context) GetPortalURL(user string) (string, error) {

	logging.Debugf("Retrieving portal URL for billing user %s", user)

	body := url.Values{}
	body.Add("customer", user)
	body.Add("return_url", ctx.CallbackURL)

	ret, _ := http.Request{
		URL:            StripeURL + "/v1/billing_portal/sessions",
		ContentType:    "application/x-www-form-urlencoded",
		Authentication: http.Authentication{Scheme: "Basic", Username: ctx.APIKey, Password: ""},
		Body:           []byte(body.Encode()),
	}.Post()

	var response APIResponsePortalURL

	err := json.Unmarshal(ret.Body, &response)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response from billing engine: %s", ret.Body)
	}

	if response.Error.Message != "" {
		return "", errors.New(response.Error.Message)
	}

	if response.URL == "" {
		return "", fmt.Errorf("the response did not include the billing URL: %s", ret.Body)
	}

	logging.Debugf("Finished retrieving portal URL for billing user %s", user)
	return response.URL, nil

}

func (ctx *Context) GetSubscriptionByCustomer(id string) (*APIResponseSubscription, error) {

	logging.Debugf("Retrieving subscriptions for billing user %s", id)
	ret, _ := http.Request{
		URL:            StripeURL + "/v1/subscriptions?customer=" + id,
		Authentication: http.Authentication{Scheme: "Basic", Username: ctx.APIKey, Password: ""},
	}.Get()

	var response struct {
		Object string `json:"object"`
		Error  `json:"error"`
		Data   []APIResponseSubscription
	}

	err := json.Unmarshal(ret.Body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response from billing engine: %s", ret.Body)
	}

	if response.Error.Message != "" {
		return nil, errors.New(response.Error.Message)
	}

	if len(response.Data) == 0 {
		return &APIResponseSubscription{}, nil
	}

	logging.Debugf("Finished retrieving subscriptions for billing user %s", id)
	return &response.Data[0], nil
}

func (ctx *Context) GetSubscriptionWithPlan(id string) (*Subscription, error) {

	// Get the Subscription for customer using their billing_id.  Abort if error encountered
	subscription, err := ctx.GetSubscriptionByCustomer(id)
	if err != nil {
		return nil, err
	}

	ret := Subscription{
		Name:     subscription.Plan.Product,
		Price:    subscription.Plan.Amount,
		Cycle:    subscription.Plan.Interval,
		Renew:    time.Unix(subscription.CurrentPeriodEnd, 0),
		Currency: subscription.Plan.Currency,
	}

	if ret.Name == "" {
		return &ret, nil
	}

	product, err := ctx.GetProduct(subscription.Plan.Product)
	if err != nil {
		logging.Warningf("failed to retrieve product '%s': %s", subscription.Plan.Product, err.Error())
		return &ret, err
	}

	ret.Name = product

	invoices, err := ctx.GetInvoices(subscription.ID)
	if err != nil {
		logging.Warningf("failed to retrieve invoices for subscription '%s': %s", subscription.ID, err.Error())
		return &ret, err
	}

	ret.Invoices = invoices

	return &ret, nil
}

func (ctx *Context) GetInvoices(sub string) ([]Invoice, error) {

	logging.Debugf("Retrieving invoices for billing user %s", sub)
	ret, _ := http.Request{
		URL:            StripeURL + "/v1/invoices?limit=12&subscription=" + sub,
		Authentication: http.Authentication{Scheme: "Basic", Username: ctx.APIKey, Password: ""},
	}.Get()

	var response struct {
		Object string `json:"object"`
		Error  `json:"error"`
		Data   []APIResponseInvoice
	}

	err := json.Unmarshal(ret.Body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response from billing engine: %s", ret.Body)
	}

	if response.Error.Message != "" {
		return nil, errors.New(response.Error.Message)
	}

	invoices := make([]Invoice, len(response.Data))

	if len(response.Data) == 0 {
		return invoices, nil
	}

	for i, inv := range response.Data {
		invoices[i] = Invoice{
			Date:  time.Unix(int64(inv.Created), 0),
			Price: inv.AmountDue,
		}
	}

	logging.Debugf("Finished retrieving invoices for subscription %s", sub)
	return invoices, nil
}

func (ctx *Context) GetProduct(id string) (string, error) {

	logging.Debugf("Retrieving product %s", id)

	logging.Debugf("Trying to retrieve %s from redis...", id)
	name, err := ctx.Cache.Get(id).Result()
	if err == nil && name != "" {
		logging.Debugf("Returning cached product name '%s' for product ID '%s'", name, id)
		return name, nil
	}

	logging.Debugf("No cache for product %s, fetching from Stripe...", id)
	ret, _ := http.Request{
		URL:            StripeURL + "/v1/products/" + id,
		Authentication: http.Authentication{Scheme: "Basic", Username: ctx.APIKey, Password: ""},
	}.Get()

	var response APIResponseProduct

	err = json.Unmarshal(ret.Body, &response)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response from billing engine: %s", ret.Body)
	}

	if response.Error.Message != "" {
		return "", errors.New(response.Error.Message)
	}

	logging.Debugf("Finished retrieving product %s", id)

	logging.Debugf("Saving product %s to cache", id)
	_, err = ctx.Cache.Set(id, response.Name, time.Hour).Result()
	if err != nil {
		logging.Warningf("Failed to save product '%s' to Redis: %s", id, err.Error())
	}

	return response.Name, nil
}

func (ctx *Context) CreateUser(p *Profile) (string, error) {

	if ctx.APIKey == "" {
		return "", fmt.Errorf("API key is required")
	}

	body := url.Values{}
	body.Add("email", p.Email)
	body.Add("name", p.Name)
	body.Add("metadata[id]", p.ID)

	ret, _ := http.Request{
		URL:            StripeURL + "/v1/customers",
		Authentication: http.Authentication{Username: ctx.APIKey, Password: "", Scheme: "basic"},
		ContentType:    "application/x-www-form-urlencoded",
		Body:           []byte(body.Encode()),
	}.Post()

	var response APIResponseCreateUser

	err := json.Unmarshal(ret.Body, &response)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response from billing engine: %s", ret.Body)
	}

	if response.Error.Message != "" {
		return "", errors.New(response.Error.Message)
	}

	if response.ID == "" {
		return "", fmt.Errorf("the response did not include the billing ID: %s", ret.Body)
	}

	return response.ID, nil
}

func New(key string, cb string, cache *redis.Client) (*Context, error) {
	if key == "" {
		return nil, fmt.Errorf("the API key is required")
	}

	if cb == "" {
		return nil, fmt.Errorf("the callback URL is required")
	}

	if cache == nil {
		return nil, fmt.Errorf("redis cache client handle is required")
	}

	return &Context{
		APIKey:      key,
		CallbackURL: cb,
		Cache:       cache,
	}, nil
}
