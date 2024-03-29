package jwt

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

const (
	// signed tokens
	JWS_X5C_A = "eyJ4NWMiOlsiTUlJSFRUQ0NCaldnQXdJQkFnSVFjRHRQY1VzUTBHbDJ5NUxCVENvMFpqQU5CZ2txaGtpRzl3MEJBUXNGQURDQnVqRUxNQWtHQTFVRUJoTUNWVk14RmpBVUJnTlZCQW9URFVWdWRISjFjM1FzSUVsdVl5NHhLREFtQmdOVkJBc1RIMU5sWlNCM2QzY3VaVzUwY25WemRDNXVaWFF2YkdWbllXd3RkR1Z5YlhNeE9UQTNCZ05WQkFzVE1DaGpLU0F5TURFMElFVnVkSEoxYzNRc0lFbHVZeTRnTFNCbWIzSWdZWFYwYUc5eWFYcGxaQ0IxYzJVZ2IyNXNlVEV1TUN3R0ExVUVBeE1sUlc1MGNuVnpkQ0JEWlhKMGFXWnBZMkYwYVc5dUlFRjFkR2h2Y21sMGVTQXRJRXd4VFRBZUZ3MHlNREE0TVRReE5qRTVNVFphRncweU1UQTRNVE14TmpFNU1UWmFNSUhuTVFzd0NRWURWUVFHRXdKVlV6RVJNQThHQTFVRUNCTUlTV3hzYVc1dmFYTXhFREFPQmdOVkJBY1RCME5vYVdOaFoyOHhFekFSQmdzckJnRUVBWUkzUEFJQkF4TUNWVk14R1RBWEJnc3JCZ0VFQVlJM1BBSUJBaE1JUkdWc1lYZGhjbVV4SkRBaUJnTlZCQW9URzBKaGJtc2diMllnUVcxbGNtbGpZU0JEYjNKd2IzSmhkR2x2YmpFZE1Cc0dBMVVFRHhNVVVISnBkbUYwWlNCUGNtZGhibWw2WVhScGIyNHhFREFPQmdOVkJBVVRCekk1TWpjME5ESXhMREFxQmdOVkJBTVRJMEYxZEdob2RXSTBOakkxTWkxUVVrOUVMbUpoYm10dlptRnRaWEpwWTJFdVkyOXRNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQTdtdU9jOGtBMlBvZy84SVZValJYSWRHWFJvVTk1RVlMNjVrYjNjNCtjOGVsQURCdnA1ZG85bk90WUhVQ29uc2daclp3UFF1U2dadHVYMlpYQWh6V3BTcHVweGkvL2RTQ0JGWDYySHd3enZvRkJTckxLOHRERFFZZVlqUnU2S0p2Sy9TWUJpLzlQVEVTNUcyRldaakM2NENhdHVFU0x2dnRHSWJPQTFaNFF2aUp3Vzg3dk12ZVB5azRDSGJWaVRXNlcrVXFMUDVRbUMzK012eFRZTzN6NS9ORnpVWnJCSGZ5aVFHVCtZUGY4aHpTbVBPQ2VvUkswNWFpVWN1UHZhQmVSdnNObUFSQzM0RDFuMHhmS1RFeEhyc3l3RklMbmI1NUpZcFRNU0VEZVJvNjdqSVQwelBnbFlTemI5TmdTcEYxRnB6dEVOdHFHN2pwRU1FSGRObVVWd0lEQVFBQm80SURIakNDQXhvd0RnWURWUjBQQVFIL0JBUURBZ1dnTUIwR0ExVWRKUVFXTUJRR0NDc0dBUVVGQndNQkJnZ3JCZ0VGQlFjREFqQU1CZ05WSFJNQkFmOEVBakFBTUIwR0ExVWREZ1FXQkJTMEx0TW9TQkdWWDRKKzBXb1lLUVB4ajEvV1VEQWZCZ05WSFNNRUdEQVdnQlREOTlDMUtqQ3RydzJSSVhBNVZOMjhpWERIT2pCb0JnZ3JCZ0VGQlFjQkFRUmNNRm93SXdZSUt3WUJCUVVITUFHR0YyaDBkSEE2THk5dlkzTndMbVZ1ZEhKMWMzUXVibVYwTURNR0NDc0dBUVVGQnpBQ2hpZG9kSFJ3T2k4dllXbGhMbVZ1ZEhKMWMzUXVibVYwTDJ3eGJTMWphR0ZwYmpJMU5pNWpaWEl3TXdZRFZSMGZCQ3d3S2pBb29DYWdKSVlpYUhSMGNEb3ZMMk55YkM1bGJuUnlkWE4wTG01bGRDOXNaWFpsYkRGdExtTnliREF1QmdOVkhSRUVKekFsZ2lOQmRYUm9hSFZpTkRZeU5USXRVRkpQUkM1aVlXNXJiMlpoYldWeWFXTmhMbU52YlRCTEJnTlZIU0FFUkRCQ01EY0dDbUNHU0FHRyttd0tBUUl3S1RBbkJnZ3JCZ0VGQlFjQ0FSWWJhSFIwY0hNNkx5OTNkM2N1Wlc1MGNuVnpkQzV1WlhRdmNuQmhNQWNHQldlQkRBRUJNSUlCZlFZS0t3WUJCQUhXZVFJRUFnU0NBVzBFZ2dGcEFXY0FkZ0JWZ2RUQ0ZwQTJBVXJxQzV0WFBGUHd3T1E0ZUhBbENCY3ZvNm9kQnhQVERBQUFBWFB0eEsvN0FBQUVBd0JITUVVQ0lRQ2crOExKT09ocjFVQitBdm5nUHdwMFYwbnRBZ295a1dJeGxZY1o4Z1dhdHdJZ2RnOUttZFR0VUlvT2Evd1RjS2ZCVlN0aTA2K2FrQ0ZFN2lXSUk2T3luaW9BZFFEdXdKWHVqWEprRDVManc3a2J4eEtqYVdvSmUwdHFHaFE0NWtleXkrM0YrUUFBQVhQdHhMQTlBQUFFQXdCR01FUUNJQ0tYdWtYUWlpdmhKei9Wc3l3WWdVQUp0MThnY28vWWVQemtxRHp1dk80NUFpQkxWei9pNTVaN1JmNGdZVkNmRExhc0piL0h2dzB5eVRyckkxTC9GekdiR2dCMkFQWmNsQy9SZHpBaUZGUVlDRENVVm83alRSTVpNNy9mREM4Z0M4eE84V1RqQUFBQmMrM0VyLzhBQUFRREFFY3dSUUlnVlo3U3B1RVltR0F6cGZ6dVcxWm9MdmxOZkcvdG4xdGxzU1BtUW9FUzVJc0NJUUN1MHNwMFdReFVVdy9tRGE2Y3gzYURjcjloQWprbWNWSEFsUE40QzBrNFlqQU5CZ2txaGtpRzl3MEJBUXNGQUFPQ0FRRUFyblEzUTZlYkJJM1RBUG5tdVNTcEY1UCtsTFFUc2xQQUM5VXF6NXRyTWlaeEdVRlVtVS9NdXUzN2hJL3U5NzNHNDBMZmpzUnZVYmxrUG93KzJldXVhVFd2VTlFOWFLYVlKY3dlOG1KL2tPS0I2cFRqdGhwU2lhQkNWODlPaDg4SVpNbEJTLzJWQ1ZKeVlVMVhCcEJVczEwYi9LNjhrQ0ZzL0ZzeUNYRllOd1BVb01ReU5aSkVYdlJTajZGRVdZYmZqRHFxV1dWY1VSTkJ2NkF1R2VFUm5FZG1HVU9Nd3FrcE9KSU8rd2Y1dlRtYUtleDh5MlZ4Q1JjcEdTL1U5TFc2MDVQQk93OCtlUGFyMzlyTDFodWhQcy9jazV6Mjd1bzFNUEd1dWVyU3JiUFdheSt3TDdNMHJWUW51SHRlemRMSHNHTWZwV0YwRUs3NTBsOWVCRGNKZVE9PSIsIk1JSUZMVENDQkJXZ0F3SUJBZ0lNWWFIbjBnQUFBQUJSMDJhbU1BMEdDU3FHU0liM0RRRUJDd1VBTUlHK01Rc3dDUVlEVlFRR0V3SlZVekVXTUJRR0ExVUVDaE1OUlc1MGNuVnpkQ3dnU1c1akxqRW9NQ1lHQTFVRUN4TWZVMlZsSUhkM2R5NWxiblJ5ZFhOMExtNWxkQzlzWldkaGJDMTBaWEp0Y3pFNU1EY0dBMVVFQ3hNd0tHTXBJREl3TURrZ1JXNTBjblZ6ZEN3Z1NXNWpMaUF0SUdadmNpQmhkWFJvYjNKcGVtVmtJSFZ6WlNCdmJteDVNVEl3TUFZRFZRUURFeWxGYm5SeWRYTjBJRkp2YjNRZ1EyVnlkR2xtYVdOaGRHbHZiaUJCZFhSb2IzSnBkSGtnTFNCSE1qQWVGdzB4TkRFeU1UVXhOVEkxTUROYUZ3MHpNREV3TVRVeE5UVTFNRE5hTUlHNk1Rc3dDUVlEVlFRR0V3SlZVekVXTUJRR0ExVUVDaE1OUlc1MGNuVnpkQ3dnU1c1akxqRW9NQ1lHQTFVRUN4TWZVMlZsSUhkM2R5NWxiblJ5ZFhOMExtNWxkQzlzWldkaGJDMTBaWEp0Y3pFNU1EY0dBMVVFQ3hNd0tHTXBJREl3TVRRZ1JXNTBjblZ6ZEN3Z1NXNWpMaUF0SUdadmNpQmhkWFJvYjNKcGVtVmtJSFZ6WlNCdmJteDVNUzR3TEFZRFZRUURFeVZGYm5SeWRYTjBJRU5sY25ScFptbGpZWFJwYjI0Z1FYVjBhRzl5YVhSNUlDMGdUREZOTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUEwSUhCT1NQQ3NkSHM5MWZkVlNRMmtTQWlTUGY4eWxJS3NLcy9NN1d3aEFmMjMwNTZzUHVZSWowQnJGYjdjVzJ5N3JtZ0QxSjNxNWlUdmpPSzY0ZGV4NnF3eW1tUFF3aHFQeUsvTXpsRzFaVHk0a3dGSXRsbmdKSHhCRW9PbTN5aXlkSnMvVHdKaEwzOWF4U2FnUjNuaW9QdllSWjFSNWdUT3cyUUZwaS9pdUluTWxPWm1jUDdsaHcxOTJMdGpMMUpjZEpEUTZHaDR5RXFJM0NvZFQyeWJFWUdZVzhZWitRcGZySTh3Y1ZmQ1I1dVJFN3NJWmxZRlVqMFZVZ3F0elMwQmVOOFNZd0FXTjQ2bHN3NTNHRXpWYzRxTGovUm1XTG9xdVkwZGpHcXIza3BsbmpMZ1JTdmFkcjdCTGxaZzBTcUNVKzAxQ3dCblp1VU1Xc3RvYy9CNVFJREFRQUJvNElCS3pDQ0FTY3dEZ1lEVlIwUEFRSC9CQVFEQWdFR01CMEdBMVVkSlFRV01CUUdDQ3NHQVFVRkJ3TUNCZ2dyQmdFRkJRY0RBVEFTQmdOVkhSTUJBZjhFQ0RBR0FRSC9BZ0VBTURNR0NDc0dBUVVGQndFQkJDY3dKVEFqQmdnckJnRUZCUWN3QVlZWGFIUjBjRG92TDI5amMzQXVaVzUwY25WemRDNXVaWFF3TUFZRFZSMGZCQ2t3SnpBbG9DT2dJWVlmYUhSMGNEb3ZMMk55YkM1bGJuUnlkWE4wTG01bGRDOW5NbU5oTG1OeWJEQTdCZ05WSFNBRU5EQXlNREFHQkZVZElBQXdLREFtQmdnckJnRUZCUWNDQVJZYWFIUjBjRG92TDNkM2R5NWxiblJ5ZFhOMExtNWxkQzl5Y0dFd0hRWURWUjBPQkJZRUZNUDMwTFVxTUsydkRaRWhjRGxVM2J5SmNNYzZNQjhHQTFVZEl3UVlNQmFBRkdweUpuclFIdTk5NXp0cFVkUnNqWitRRW1hck1BMEdDU3FHU0liM0RRRUJDd1VBQTRJQkFRQzBoOGVFSWhvcHdLUjQ3UFZQRzdTRWwyOTM3dFRQV2Erb1E1WXZIVmplcHZNVld5N1pRNXhNUXJrWEZ4R3R0TEZCeDJZTUlvWUZwN1FpKzhWb2FJcUlNdGh4MWhHT2psSitRZ2xkMmRuQURpenZSR3NmMnlTODlieXhxc0dLNVdiYjBDVHozNG1taS81ZTBGQzZtM1VBeVFoS1MzUS9XRk92OXJpaGJJU1lKbno4L0RWUlpaZ2VPMngyOEprUHhMa0oxWVhZSktkL0tzTGFrMHRrdUhCOFZDblRnbFRWejZXVXd6T2VUVFJuNERoMlpnQ04wQy9HcXdtcWN2ck9MeldKL01EdEJnTzMzNHdsVi9INzd5aUkyWUlvd0FRUGxJRnBJK0NSS01WZTFRelgxQ0E3NzhuNHdJK25RYzFYUkc1c1oyTCtoTi9uWU5qdnY5UWlIZzNuIl0sInR5cCI6IkpXVCIsImFsZyI6IlJTMjU2In0.eyJpYXQiOjE2MDY5Mzc4MTYsImF1ZCI6ImFwaS45c3Bva2VzLmlvIiwiZXhwIjoxNjA2OTM4MDI2LCJqdGkiOiJseHg1R1FZNWxURXhjaXZadlhDdjN3IiwiaXNzIjoiaHR0cHM6Ly9zZWN1cmUuYmFua29mYW1lcmljYS5jb20iLCJuYmYiOjE2MDY5Mzc4MTYsInN1YiI6ImYyYWMwNThlZmI3M2JhMTg3YWE0NGVlODM4MWU4Zjk2YTBhNzJmMjZjNDBlYzMxMjFkYzJjOTA3ZDYxNjZmZWM0NDljM2FhODk2MjY5N2I3YTA4M2IxMjY5OWRmMGE0Y2U3NjMyZjBlNDFkYjkyY2NlOTQzODk2NmZjZDNkMzAyIiwicGVybXMiOiJjZjJkODNlYS0zMDZiLTQ0YWEtYTNjOS1lODFhMmM1ODE2ZTgiLCJqc3VybCI6Imh0dHBzOi8vc2VjdXJlLmJhbmtvZmFtZXJpY2EuY29tL3NwYS93aWRnZXRzL2xvYWRlci8xMC4wLjAtOXMtbG9hZGVyL2luZGV4LmpzIn0.iULNQ5nBAMNapGlDQ7XSFhFVxdsKhLSbk4C8SCQ9QpXLklFU7ONXq09eB5pOjmlX-Huxo_mIybGzaGlszsT5JU5xmYuXsUDIPBhK-nlkw80n8U4f1IdUgjAyVMqS-YvgGWVFtBPqyA4EQB2ZpO_XX8sX5SVm6hx4meFyl8tfYkbXdLi1sFvPucQIHchakaYDLxSKPP7Y5as06r2WJWUc7RXn302MeBoPdFHEjFM2hBhy3stSmloMuJFvHMDTKjm7Yd_P7eZBjMwxFvx5GnEc2iKI6SJ1v-jq-jW2VDJ87ACz3tJUv4lgbs6Vd0iSquvJ0RECOywJFMT8sHM3Zxw07A"
	JWS_X5C_B = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsIng1YyI6WyJcbk1JSUZ1ekNDQTZPZ0F3SUJBZ0lVQXVEZWx1UFE0cjlYRi9PVXprOEczZ29BVk4wd0RRWUpLb1pJaHZjTkFRRUxcbkJRQXdiVEVMTUFrR0ExVUVCaE1DVGxveEVUQVBCZ05WQkFnTUNFRjFZMnRzWVc1a01SRXdEd1lEVlFRSERBaEJcbmRXTnJiR0Z1WkRFUU1BNEdBMVVFQ2d3SE9WTndiMnRsY3pFTU1Bb0dBMVVFQ3d3RFJHVjJNUmd3RmdZRFZRUURcbkRBOTNkM2N1T1hOd2IydGxjeTVqYjIwd0hoY05Nakl3TVRFNU1qTTFORFU0V2hjTk16SXdNVEUzTWpNMU5EVTRcbldqQnRNUXN3Q1FZRFZRUUdFd0pPV2pFUk1BOEdBMVVFQ0F3SVFYVmphMnhoYm1ReEVUQVBCZ05WQkFjTUNFRjFcblkydHNZVzVrTVJBd0RnWURWUVFLREFjNVUzQnZhMlZ6TVF3d0NnWURWUVFMREFORVpYWXhHREFXQmdOVkJBTU1cbkQzZDNkeTQ1YzNCdmEyVnpMbU52YlRDQ0FpSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnSVBBRENDQWdvQ2dnSUJcbkFNTC9HZnR2bEN6eGN1eUZ4ZFRIbFdmY0xMRTBUU0p4ZWN5ZVhZZnNLWU1reHR5Rm1kc2UyUDF4VmFiNUNmUnBcbjNTVmJpbG94YjYwdG9Vbm1LTDIzdHNUN3VuYnlaZGorRHRsSGlZTnU4ZmlaRXNra2UrSTNXQnBKNmc2dFZUVHhcbmRjZTdGZmsvb08zRS9ZY2lEdkVFckZRMXRick1KMkE0SDJhU2NONjJiU01naGVEbG1jclUwTjliWVZXS1JtbGxcbjZBMXMrcjBwaTY2UEdMYXdYeUhtWERzOGgwRE9XcC81cDdlRTBtTEI2bDk0WXJiWHpzaG5KSHdZVy9zN01raE5cbjJZWTBmWFZzTDc0TUt0N0VZMXA2WktleVcxc0hUSnpIRThXa01QQkQ1cVBHS3hlcThsUW5rVjFYYmlnV0dQOGlcbjI2UGF4UndiaTRkYmxkZmFRYzgvQ05WYXkyUmRQUER0SXN2YnZnNTQ3My90aHA0WjcvVm5PK0t6cHZxem9seHRcbk5jcTlDYlY0TWtFNzFpdDR6ZW8rOFd1V1dQbVh4WGpvNE03LzZPak9DU0prSlpzaDc4ZnVzWEo1T3ZkeWpIODdcblFyLys0RG1tVnYxR3JnZmhXc0NDZG5QczJhR0h3Q3ZIcDYvUDVpUnhSV2RLSVNIU3BsdGZod1hrSUNiSzJyMnFcbnVnTU10KzRibkdLaFRQelNvRmx0cVBQMWZwTGQydFJ5MURBWDlMOUxBUGJGdnVvYUQ5dG1pUUExNXNVVnNvckFcbk1yZGpkUkJCVTdYVTNLcEcydm9yRmNzK2cyYnI0bGJ2UHovS0M0dmtyNzAwUzJHa0NFb1hlRUNUaCsxSVdxNnRcbnhYc0ZsWVBLNlBrREpuWDNVRkhoN2NDN2xnYWpwbTF5SHUwS0RYTjV2QmlqQWdNQkFBR2pVekJSTUIwR0ExVWRcbkRnUVdCQlEwOFk0WDFqbnkxRW51R3Q3c1pUL1IwRzJna1RBZkJnTlZIU01FR0RBV2dCUTA4WTRYMWpueTFFbnVcbkd0N3NaVC9SMEcyZ2tUQVBCZ05WSFJNQkFmOEVCVEFEQVFIL01BMEdDU3FHU0liM0RRRUJDd1VBQTRJQ0FRQmdcbkNkQ05JYWJwWGYvVGpTQldFSzA2b0I0L3lzVXArV3B6ejZTKzJ3dXNjMXJoMXpYaDJEcW5VOVZyaFBVaURGZWNcbkQ3ZXVlaDhPcHI0VlVZbldCbHRJQkMrVDFtKzV2RExXN3RZZ2R3SCtBYVpzNGx3Si9qcFpSZFI3dUZEcndsd0pcbjhLbzRHT3VLTWw0UjdaYXBGSzZuY2NWUEMxaE9nUXVtWG42NFR2WTFHZ2lZRmczQndEck9NYnd2YWNRMFhKVkFcblhiQW9IeXZLdWNhNzY1c1BVbGFRVURLWHp6UXZFMm41eEdPNzBGbEVsakxpV09aMVBLdWQ3dUZOWWgySlQyWWVcbnMvbWpORWttVkpsQU5JTVJDakJqdmZDTENtRVRRL3JqWVJqRTBQRzZtNG9MNDRIc01jdjFxaFZrVTNES2xOTWlcbnVGU3ZIdGVvM1cwQkFLQ0NIWGFITTFZVUhOWU5kT21qSDd3VkpXL2lyTXFiMHhkM3F3VHhKRzdwNHFBSnMzdXBcbnd4aVNxZTJUajB4bXV5emVLZURpdS9RUkc2Y3FRNnlkdGp0WC9yQ1VoS3gvVXgvRCtYWlNnY08xcm4vNW41aWdcbkZRdm15eU5tNE1BcGNvd2srMEQ3MW0yUlBwQWdQUjIxdUhsTmZibytBcW5NYUdSNzBFaFVWSVBLalV6ZjlWRWFcbnVONWt4bi91MkRnUGhHNUVvQkprUWs0WUZqeml1NS8wclFoUzlickFsS1JwYTQzTmh2QlF1dW56REJRZHR4VlVcbk90VlhIelQxaDB2VEpVTlh2MXhJVWZBeDl0cXFLTXRQbXl5YVpoUWtlcDdqV2llYzJaQ2Z4SjZlUm1FbzQ2ZUlcbjFubndjaEJEWWlVSXZKZEkzVnNzZEVRWFlrN0lvOFpURFN1S2REUXAwQT09XG4iXX0.eyJhdWQiOiI5c3Bva2VzIiwiZXhwIjoxNjQyNjQwMjY4LCJpYXQiOjE2NDI2MzcyNjksImlzcyI6IjlTcG9rZXMiLCJqdGkiOiI5ZTQ1YzZlOC1iYTg0LTRiOTUtOWZiMC05ZTAxNDgyNTdiMGEiLCJuYmYiOjE2NDI2MzcyNjksInN1YiI6IjEyMzQ1Njc4OSJ9.ddocoK-CWBnfUltmx2HGmEk3i5VwV03oxx7YloMCZlJ9di8wVSIPftZGzdQdOT14BHumfqzyIQdfu2WYuhDQ0UCqdvfVQVPGOJuW4kYitakKJUMAgyzePWqw3gU_L51VepelDLN90ZV4F3qzLz_6JTkL_Xo71xXc4hLUJG91yTLYmzoyla7EFZhbkDwkDPGlLy2AQGEwWgdUVeICfU6SFCfbfOhES2dUGdVvjTFzX6hXWBKj_3V2xthHUyT6R3BjbEAgxbUqlaV3sS00b3IFq4LfNlm7DrqJI2DEyq1yoUh7vD85S-Csdr0ngoxWj9aR4OQVGvLHtvNmJZdK1FvypEf2vqcCbmXnAAfJZk08b0aqUA4Ns8Gmz6o9N60rc4LVBAzpOMaSrc2ih3PhRJBBvQeg2_FKifHx9s8-qeypmwYTpVIBZFgR2cEyhx-QyhZzy_YGXOFJ0XhpAkXAIW3L326LuwGZYfTgbnnxxMkg8VlgemyW83HKX-A9sHGjXpO6kJrRDNNPvsMvgih4_c3CizCoC2CRfeZ945Uc9dQf3BucV1PAzV1kq8JzNdoqj28RsUX5SQhSZ9ProQ9XyPLQF228lO5CdIxQCd_RLx2eqWj_xcwsLgaLolnTzfsuVqc7tBIYXFs57bT4UraMVam3wQOw6qjXKVDdnIoDl7tlM_g"
	JWS_X5C_C = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsIng1YyI6WyJcbk1JSUZ1ekNDQTZPZ0F3SUJBZ0lVQXVEZWx1UFE0cjlYRi9PVXprOEczZ29BVk4wd0RRWUpLb1pJaHZjTkFRRUxcbkJRQXdiVEVMTUFrR0ExVUVCaE1DVGxveEVUQVBCZ05WQkFnTUNFRjFZMnRzWVc1a01SRXdEd1lEVlFRSERBaEJcbmRXTnJiR0Z1WkRFUU1BNEdBMVVFQ2d3SE9WTndiMnRsY3pFTU1Bb0dBMVVFQ3d3RFJHVjJNUmd3RmdZRFZRUURcbkRBOTNkM2N1T1hOd2IydGxjeTVqYjIwd0hoY05Nakl3TVRFNU1qTTFORFU0V2hjTk16SXdNVEUzTWpNMU5EVTRcbldqQnRNUXN3Q1FZRFZRUUdFd0pPV2pFUk1BOEdBMVVFQ0F3SVFYVmphMnhoYm1ReEVUQVBCZ05WQkFjTUNFRjFcblkydHNZVzVrTVJBd0RnWURWUVFLREFjNVUzQnZhMlZ6TVF3d0NnWURWUVFMREFORVpYWXhHREFXQmdOVkJBTU1cbkQzZDNkeTQ1YzNCdmEyVnpMbU52YlRDQ0FpSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnSVBBRENDQWdvQ2dnSUJcbkFNTC9HZnR2bEN6eGN1eUZ4ZFRIbFdmY0xMRTBUU0p4ZWN5ZVhZZnNLWU1reHR5Rm1kc2UyUDF4VmFiNUNmUnBcbjNTVmJpbG94YjYwdG9Vbm1LTDIzdHNUN3VuYnlaZGorRHRsSGlZTnU4ZmlaRXNra2UrSTNXQnBKNmc2dFZUVHhcbmRjZTdGZmsvb08zRS9ZY2lEdkVFckZRMXRick1KMkE0SDJhU2NONjJiU01naGVEbG1jclUwTjliWVZXS1JtbGxcbjZBMXMrcjBwaTY2UEdMYXdYeUhtWERzOGgwRE9XcC81cDdlRTBtTEI2bDk0WXJiWHpzaG5KSHdZVy9zN01raE5cbjJZWTBmWFZzTDc0TUt0N0VZMXA2WktleVcxc0hUSnpIRThXa01QQkQ1cVBHS3hlcThsUW5rVjFYYmlnV0dQOGlcbjI2UGF4UndiaTRkYmxkZmFRYzgvQ05WYXkyUmRQUER0SXN2YnZnNTQ3My90aHA0WjcvVm5PK0t6cHZxem9seHRcbk5jcTlDYlY0TWtFNzFpdDR6ZW8rOFd1V1dQbVh4WGpvNE03LzZPak9DU0prSlpzaDc4ZnVzWEo1T3ZkeWpIODdcblFyLys0RG1tVnYxR3JnZmhXc0NDZG5QczJhR0h3Q3ZIcDYvUDVpUnhSV2RLSVNIU3BsdGZod1hrSUNiSzJyMnFcbnVnTU10KzRibkdLaFRQelNvRmx0cVBQMWZwTGQydFJ5MURBWDlMOUxBUGJGdnVvYUQ5dG1pUUExNXNVVnNvckFcbk1yZGpkUkJCVTdYVTNLcEcydm9yRmNzK2cyYnI0bGJ2UHovS0M0dmtyNzAwUzJHa0NFb1hlRUNUaCsxSVdxNnRcbnhYc0ZsWVBLNlBrREpuWDNVRkhoN2NDN2xnYWpwbTF5SHUwS0RYTjV2QmlqQWdNQkFBR2pVekJSTUIwR0ExVWRcbkRnUVdCQlEwOFk0WDFqbnkxRW51R3Q3c1pUL1IwRzJna1RBZkJnTlZIU01FR0RBV2dCUTA4WTRYMWpueTFFbnVcbkd0N3NaVC9SMEcyZ2tUQVBCZ05WSFJNQkFmOEVCVEFEQVFIL01BMEdDU3FHU0liM0RRRUJDd1VBQTRJQ0FRQmdcbkNkQ05JYWJwWGYvVGpTQldFSzA2b0I0L3lzVXArV3B6ejZTKzJ3dXNjMXJoMXpYaDJEcW5VOVZyaFBVaURGZWNcbkQ3ZXVlaDhPcHI0VlVZbldCbHRJQkMrVDFtKzV2RExXN3RZZ2R3SCtBYVpzNGx3Si9qcFpSZFI3dUZEcndsd0pcbjhLbzRHT3VLTWw0UjdaYXBGSzZuY2NWUEMxaE9nUXVtWG42NFR2WTFHZ2lZRmczQndEck9NYnd2YWNRMFhKVkFcblhiQW9IeXZLdWNhNzY1c1BVbGFRVURLWHp6UXZFMm41eEdPNzBGbEVsakxpV09aMVBLdWQ3dUZOWWgySlQyWWVcbnMvbWpORWttVkpsQU5JTVJDakJqdmZDTENtRVRRL3JqWVJqRTBQRzZtNG9MNDRIc01jdjFxaFZrVTNES2xOTWlcbnVGU3ZIdGVvM1cwQkFLQ0NIWGFITTFZVUhOWU5kT21qSDd3VkpXL2lyTXFiMHhkM3F3VHhKRzdwNHFBSnMzdXBcbnd4aVNxZTJUajB4bXV5emVLZURpdS9RUkc2Y3FRNnlkdGp0WC9yQ1VoS3gvVXgvRCtYWlNnY08xcm4vNW41aWdcbkZRdm15eU5tNE1BcGNvd2srMEQ3MW0yUlBwQWdQUjIxdUhsTmZibytBcW5NYUdSNzBFaFVWSVBLalV6ZjlWRWFcbnVONWt4bi91MkRnUGhHNUVvQkprUWs0WUZqeml1NS8wclFoUzlickFsS1JwYTQzTmh2QlF1dW56REJRZHR4VlVcbk90VlhIelQxaDB2VEpVTlh2MXhJVWZBeDl0cXFLTXRQbXl5YVpoUWtlcDdqV2llYzJaQ2Z4SjZlUm1FbzQ2ZUlcbjFubndjaEJEWWlVSXZKZEkzVnNzZEVRWFlrN0lvOFpURFN1S2REUXAwQT09XG4iXX0.eyJhdWQiOiI5c3Bva2VzIiwiZXhwIjoyNzQyNjQwMDM5LCJpYXQiOjE2NDI2NDAwMzksImlzcyI6IjlTcG9rZXMiLCJqdGkiOiIwZmJiNGZiYS02MWYyLTQwN2ItOGNiZS02YmU0Y2M5ODNjZTgiLCJuYmYiOjE2NDI2NDAwMzksInN1YiI6IjEyMzQ1Njc4OSJ9.TRQ9Of46O191zJQ3b2G6CIHhfOGUvF-gmNDOMllUdhPmFVh8zYR-gipxyhxBjMWVUK5WcoAkpw-pbQqSHcc8pR3pPlMJr1BUJ4LP1p2oyciRlbeQcxEhSTIu2svlyq600iWj0b_GCqowEVhS29uDvFrdaZgKDWOY2Pb4SRdHstTd5EdIV6pDBM77Ph_KNOEKrb4mkk-5l6nnnZo0XKu240ZkGX66PxWFK5zPC_KhoBnwxpQ5Np0oSNZx_KdP6mWwse6sywuJxjLf4Vkk3-_xT87ZaIls-WxprWCSNUfLRTT19oyEQ8lMFJ-SrJMbA4hacc2XwLDHGtvzgazXkK6h7-g0FeUCufEMcvti8XLkH9N0Ds-Lz6QMOU-17QcJKofDuw5Tq7LvXzz6lDGhn3FSK8L0rY55YA1ckyX_G_47hcQegWmSUbQocuZNls0OiZNbCPL1js1WT9RYweLHw-_k4NSmAcn__cysoj3N7psstxWs_-bBBfyG3BSrv5MLbTWvsuJNmTX99WQmgt6dCvbTm0gps6SW_B8YrloEB_dT_tIumV8gjewLDUNs6df5OxDRj5CalY3Py2XoFZbJ0DnwXUY-llLvQvs-ec_OaOlq72P90JxXwbO-EQ2dA7l5dTCpht-EuebjamyweyIrL4ZvGKWA7nAcaUDcSmwo0ptFWwI"
	JWS_KID   = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjljY2E5YjQ4ZTkxYTNmYWUzYzJlY2JlN2RjYTJlYTNmMTFkODRlZmEifQ.eyJpc3MiOiJodHRwczovL25zcC1kZXYuOXNwb2tlcy5pby9kZXgiLCJzdWIiOiIxZDZjZGYyYy03NTgxLTQ5OTItODZjMy04OGIyMGNkMGE2MDMiLCJhdWQiOiIiLCJleHAiOjE2NDI3ODkwODAsImlhdCI6MTY0MjcwMjY4MCwiYXRfaGFzaCI6IjZPTFM2VVRSWFJRNzZhNVNLX1VhOVEiLCJlbWFpbCI6ImFkcmlhbi5pb25pY2FAOXNwb2tlcy5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwibmFtZSI6IjhjZTE1YmI1LWJkMTItNGE0Zi1iZGFmLWJhY2UwNGEwN2FhMCIsImdpdmVuX25hbWUiOiJBZHJpYW4iLCJmYW1pbHlfbmFtZSI6IklvbmljYSIsImV4dHJhcyI6e30sInBlcm1zIjoiOGNlMTViYjUtYmQxMi00YTRmLWJkYWYtYmFjZTA0YTA3YWEwIiwibG9naW5fY291bnQiOjQ3Mn0.guwxQhzb91brzBXLrl2mTT-j4NWAFELkyMhrz2MwG44n0E7-ZQ4k3s_nOhhu2GYiCVEZbmsewxdX11UE5rkv8xC33chvhPHhfHs8AOZszYz6KdasnEqWzL3D_H0c_687RpVM11KSO20ArGrlUd1cLlaQmojel1P1QdZ4uCEd9tchjm4cjLFSxGRwZmffJpKn9-Zlsp8nhEnbPmtar_AG0B46Rw7uzHVMnC5aid5LjIpSDg2a9qbKAcVvD3kTQWAB9mGp4PZfc5GGqaLsZYKrt9w9_uAsB4QjWhNg0Kxbci6CmLRVBrT_li3IhXBc1ruDiv7EtKK0kii_0fgHlafO7Q"

	// encrypted tokens
	JWE = "eyJhbGciOiJSU0EtT0FFUC0yNTYiLCJlbmMiOiJBMjU2R0NNIiwia2lkIjoibmluZXNwb2tlc19wdWJsaWNLZXkifQ.Np1VYSJ-78j0AnAnPx-ycfXPVfQGu429TX2rlp1YTYyc5dhsZcdxv_yi5waYISvo-07-9j12oXz8Cnx3LXtbZ7cahp44GVT-0Qq1lXcRdMkNnp3488PTMo019gqH_1TgNP9kCpUXse6ErYm_iShkD2BF_Fm5yP5-_RRIaBFx4VWgFWPmj0H3Ub2_3bJqi3HHPadb86W148NViMPitMk_8kDtQomFBXuua1IHDuv1gXiQtWEkvVDLPGVipvRscuD0QJbanZb7kQgnVgLebXrTOlhCzQ-Po9m9HhvjaOvN9FEb9bu7KmbXNw_1yhWk9s47uK6Q-oYg7AYE3odwmnQXGA.clVEenkzTzh4UGxz.5IYvp-5LFX9zXTMDbxrRvehE0JeuD6IPpxJD9cbIGB8VcZTNVNC9__Fhp6KUBuMAJajJAbG-4jFSL9vgsYvDpFP9FI5gDzCqxTDQQDEpIS6d6Ej-VbEsRraty4mrr9nZaiBT3VP8NKOGqN5cS7X-upp_1zFXK9xMJY9CmMfd9iVbPBaRJAKRdFi75xvhYG9Pk8jI34nisaat-ZHu6EBvP63Q8FLQS6mr0sRX4IewPH6RHrhyhCcGZkhh6tbAgYWXl0wBLg7oU4TaM56bLbKdfTMn4eQF0RYRdAxivIpajXqFmvc00zf8I3whvs487HJq8e4vlLYp1UtGAQcIWdCu0SQ_f8LD4s__vhx3tepUwMRgl7BM9bxmfQoMczxqY3ai2y0crc_jO88hjjcbbG7XLHsu8wzh7kIpiSZ1Vn4n8QR1j7DXzeIC6r517Tt22LTcOJ6lTSf_-nvtWH2YezGQLNOYKlig-COKyHBOxIUPeD1TVpdz07R-euD5UqhhBWsj86RkZr2uPOerxDvCMbiC84VfFpJ2aRpY9_KbT0Wk51SNELkuMzJqpyJ9vUFvaFDDCQ93jxAxS-37Wfhx7tYh1UwQR3vSkRWFW8dcdbeh6MAMhLGu8kOIakfR6TZKpI5_YPdjkAQAw9Hfke9hZKzk7yUULzsryCRzxXwpz9bjq77zHtjHOfWUHCVgOjJXUcMSPOTcMve7wA_NByzIaUfIw6nsIWQBuKNS47rd9hSiszRVOv-sPaUB7uXD3oMeRT-VY_zyDK6ke3udo6CmPgxL9PrRJrrLKdzwYfKWZaO7OHiEygo9hvif5BIvyA8wpi8x9hNscGzSXFRk8ScpSsjqflKyRfsALPCuGocm2V_MvrjhncwBjTffu39t7jSh72vuBwOnRNeBO33LZuuPq1VMog8KJP-VkJgZ7ZXbKYjsoBfDGX_OSZipmLuFUnYYzWrgbo0NTmP2Tnfwsxa9T45Gkm6CdSGigBw0XhpsuvoGHhQJrKVzIAZHre6qAf93giexihxYlcqY3Ax_8zFCh02p1TYVNEH0nlRYBVIQprD3SYYA8bgr4XZPHPb_JTu-3V8wfEWb3wtcjAQp2jxgYrWQi_GJCoOGVHxGq090uSzGx5I6oyjHRPBJm5rFffeCmuDIgQZy9kSPIFSmzG2tcsSy3EV6-gxxa_PpXI980PZOQ73UdaxxnZSxv6OcYFStnL71CYzO8BaHMPAFnsipXBxbeRWugfhRLfHXCGO-PYo40t8tJoVK4IyMoGv0YEBr62klLd2kbxjXPQgmek4EALmj5fNk3JCTtCJNid-lVrsRn7W2FI0h44IJRGXJiUBU1xmG07_1l1GB8hK8PJTBTfE-KdgDy8vA1xYB2HU2_4SyqtphTo_J2ydK1J9sr4lvo0fsC1MExZC23306ZHYy7LPeOVTEm-JYqFwe45_SRENTRI7SQDm1eTy2jCk3I-JgeBC8LrG38QEhBl8vbk5FW-63c_eQvQY5FWVxLKkSQxHZdmamKnqIubvh665BBV2UHCAV7wIaEkfcS36UPqs9yQt_pzXI5vn4L0jUgh8IST2idUlxVSUJMzufEnzA8BQLLP40x1w0uat6dMcOlMmQ8pE-4GkKTI8_-URf9rpq2oddRXP3tW-WtYiN6wJi53D6VFiKhBUFDyVFFIsJ0z7bk5OxqYMBG3UjBk_PlA1KmDsGFsD1raR6HELYDOmC7wAeSxFkZxR8nETaTPkR6vWKAi86ekEygmwF6nIiaBlx7RPZPkMNulbMz39ZT4EE5kqeu33rwPKZ7AfbzV9pPGGn_9iroZbAiy3Jw8ZRTG0S9HnSVUHqZmzdTZ3Ra3SHR1J_ITw3SyaDnv99KI3n9z_Q2JAYZlMpvqZG5MGjyNJzwx5xbSkTuXbrUTMgUs10avT48lDPgUf0etfzODcONaI-VAC3rhEhSHunB9iUd9DBHH9vVvVB76x3oWVA-4Cqr4DMojhCpJAq7NG8zHqKQQvJEq02nsGN05gQxjR4-7SGIVYUnaaRmUIwuZ6RO2ZI3Lhzzzvefmz1Tn-hDsqsjKNNX5PM4pcTZIDIcXW1i7K4ztsnjwE4JZQlzi8KIxQnLeHeQOrqmOSbzgtJeTuIpMVS1znbovvQgj-762SZDYYMlYHPETotjvQMY32qqzffre89ug_TMqCbE56Vx-pNxN2IvlLjwwAyc_mtyBnd-FCJ90G297Ol_R5jM0HVnqUmm_-ms1W-5IWCijlGRh3rPs4qhEFiSnmGmVlC2VmRLHcCb4VpQ9KCg9tLpsYCW2Y2c17TM0i_K9QZ_135Aieiqm4cUr7e2eV0FhfUWy_EaBh0pfczitGIN1drWHMoWz0TBUbRFq7eoTO8rpXGixEWtvO0DJk8pQLiiiRNVtG6ac1R1F7Ys0Wi-16HGngJ8lLJKOg_bM5FoRVv61ZoSYt9VuAwvmE3l1gkJ8w9_-n2v3EfSuTPipVVBzlsX5i6thYR2HYPPJUyQbfy1pkzO80gc0v9lmGvZtSECMIWuTVq1f99SOv99aoGcTu8XJBtqn-0uhMUiwgseTbv9cCyEBtedhiUpyJe7atqQRNRpuQ8bgvmNFBhkW7iFYNMOpah8jXVifyQobTET6pkq7zYteWvDTlBwAzVHXnkTq1aTiXoopobo4D1cD0Ynf2Q7kcvyqGZAXU3BcXe7Z6UUUDA2WxEz7z4LgAMyuez2xuQB6DP2FYnHUH1FMqaqSxbFLFSNOJEbeGJw18vRE0Q_gcmniPXm9NwN334aCzzE6aD9AFaPDGDiaqoiEVpi5c8WNC6PK3ONDCVkpnOsFOaE56NQxviCv7DRecJpfPsnXsyH7kyEjfYrPkdBITZDqcHaTOdKiHZkAEF8e0e-s9pSc8mN0kbijj0U15LB6jL3h6lmJEap-uUXeuyy7XX2vklaBhervyPFcSDvc-vSFH55NY1-B0YjCpYlD89OuzORMeUmlpHAk2fGfq9ziXy6GoE5UG5kkTuw6eUmm0-22e39SpiALqySjMQcre-xeCGvgJ7OVlCq5MVk0IkIEF0T-JbGeDcUtH18YYhub9rJfJQSWW1giOj4yXbwpYj7ji3WtGeOMqvB1tZtoUUYAtW48SJF8tmeoAkXdC5eO1V-j7EWlkzl9av3GHOSOuMjVnfdL1bL9ZDSby9AYnY0xQxSuWmU67PUTDwK2AB0IaUi92q1HENaOXHcRpjvPl5sBahhfKZki40e6QBUW_oOa4_hpLo92ZiVfaKSgx-tOpVrqhOx6o8tZQlZ4Ak8n28s7Gaz1uuIjpsU2WUFE3mRSxpZ1QqPabn6gURZJCCfxU3DlDVO60HoSNpWStp26qiwP61fRB8Ba06MPSXPqcuUSUlkG0Lq0809xZkdqZngten72dhlAWvXf5c5cjZqemszzqOqPHA7cPKIMuo9yvsIRYh9_zSbeAxX-kWjzirWgUkYLLsmyNUq5ByZuxqC0okEC7S1NrQg8A-R5Zlyu2vOyXhzIDH1Lx0WcnJFeE2p5MI7OKcRG8_QYA0656nQSSVDDz406yZ2P8OLeJfGsFspKg3jsg7YJ8fEYl9g-aoMVw-dN07Jx1DqySFwRf2TjhmZ6e0OUoDPeNoLQCJ5wQ8X0chFXykPhqliFQZYTPqjJRBlWxUbMqSkLAcJO0lbdcecfEZhP-UVZ4OThxHzqPPm1FfHJdWBJTiIxdzWGQyEYG5QU_3Z-Jv7mSwHXqubiG5iZsvQArbdUnred41RT-0abYrv4fQgRBu8vXkDXyvEozZGxbujpHo6iQq8to5GtNEqSWPUsghxWxvwH9V1ElMkGsTjTbaQIDDVXo_sojNt5sZsb1SPaPM51HkD4UYolRtJREE6xmo9pltKC-X-BGoigRn1HH9q7H40OLYnPpk3cqgnGNJng6jc4FSu1V1XDjq74qCiVO6X0_3m6NrNFbhG5RUZ0WVn17_LyKQxJSrQQfUeUs8ql9yZ-jGUOhDNA2oGakIXR_jq7aRCQ78iv7wBdCxvsYuktxOxRCVlkmOIKii3K-75ZV_Ro2o0R2TUBP5oGIpQ_YDtoOC8gAKGST6ufGjMw1-Mhcj50Eo-GnPjzBk7yh_KvoKIO0GDh4Uf_Xx4IvtAKpRTH59oBaSSZoaC_a2fh8YQdrpfn4xCjD0xxmoehWYuAwfp2QFdCIJB_efpbQiZpLmb54peKVtZSh79c3BNaWy6V0jD6YfCeklZigvXB2bmo69t_DGK8RNKyAXdo8p_FNV1f2inUfx2aVmRq-ySpubUFZSrGtRzt8MTO14oocFi1Er6xPH7JKaaxd4MChk655QjrRy6_3AJ0-rdBKadqB1d8xrs9a31Mu6uHWXcFIxSTUKm5HO4GDUQgwMeuGikjHqGITHQlkcZxNIzpctrZ2721gd1Ce0_OVMsk7oSqnI7AX8dRkono-HfiLD1tiekZIfnRPhGJ4q7fOck5HQrMlJ5cuqa6z2CeMUJ6nsoxnWqPNDs0kbtWeZAj0f-jsRWZTv-q05c6XtDODpXy2J4G1QFLw2BL6yyPolCB4mfPuOegFD6bNjUSsnjnzuBU4PpYyMPBPdClgBeZZ0_jI9f19FRgm9iU5Bf2YoSA2ZD89uJakbdyfwdGPbkHHn_E5W35tTcX0kUOh5ArLLUTsIOOKhGMf-e_BNxM5AT6fjoTDGi0QKXHEZZ8rCTtXi_ZkYOTFVnMY3m-n_JDL5OxVwfvm0WPYSrZ142iEHgoOAaCm2DAxd5wwW2pWDG0-UAPyo1NDT_ZI-5VJqT4NvMufEHfN3I3UoqGd4ocGeh6WUxdlYOVSi9kqM1KhcYaQnxHdWyh1Qm2KZ5dX3Xzl5Cpam6g6b5x5JRBQTfsFjsUJl8aqvCD6ehVPhexafUjePgVwRJ4t1ufIP5JzXj02RS00TE7HyHyK3gJRd54J_FI6eIpGc3jETDNeyRqipTjJhV9qaJbp8ar1pkreO6Lt_Fur4Qm-WMh9O_iuLTlNHpVhNxWEt-odXzIaYb5278XSMiKI5zNrQIJoE83WZDUJeCL8GYbye7JqGA870GGBU-rCMn7_JI4UwpD9W0pCFnTQiFuCqve8pa54onhMfEp96eNp3IsyWHtoyfb0BAXG2EBrksBux5jJrfHk41AsnEfok63KkP3F9CG5QL7V2NTJnKLbLd7CbdFOqcoUMMfJxYgeAxwLRLd2Cd9KHhZaKoXdruPT2QmgZsJkC0DY3eXJVZosaDNmwTqkgtG6gONYEdLb30sOO7vgX7xbWwsHRjnZ3Iw7pnvqcXgpeFVvnqZUBUifFK3sNNKhQMnaCy7YnaIbp6o3JjMhJjQuJTeZmOpiLCGpP_Sdkdy2XpsqqLjvDvnG4PRGkGCnXn57nT9d2qeNMnuH4qvbvHQwsxtuEMYl1rAiVQXV8DMgz0V0Eur4dmKoAbZ267UVNl5m53BQ9JRhoscFc_NR6_aekQpfr84ZT6_LMZ4XQFF5_aBm2LR4E0iWrY9Kpsuiu1KkaHWIckkReLseQRSdCYs6UYUTvmkozAiwzphzgPCX3o4Lee1U3uP0dxaKdcvnOfHTRPtzuOha2Pj9GA7sMM8Vl_ohSTZKozabrZiMkAer0Z6H_XpwRugbz2XR8bPG6-zkt6iFWoUWG1qdqDJNP1i_zPG6xWiSY_G6zbUVa1l0A40vFYh8FoEYpm73mjSMeVNMHVl5QVsmjl2K4t3h3AAhX2I6smnveoITP1LiMo0appsvspWPz4g.6_HXe0RQHU77t_mpDaR0iQ"
)

// Declaring keys and certs here to simplify testing. New ones can be generated
// by running the following commands in a terminal:
// > openssl req -x509 -newkey rsa:4096 -keyout private.pem -nodes -out cert.pem -sha256 -days 365 -subj "/C=NZ/ST=Auckland/L=Auckland/O=9Spokes/OU=Dev/CN=www.9spokes.com" -passout pass:secret
// > openssl x509 -pubkey -noout -in cert.pem  > public.pem
// > openssl rsa -des3 -in private.pem -out enc-private.pem
const (
	PRIVATE_KEY_JWS = `
-----BEGIN RSA PRIVATE KEY-----
Proc-Type: 4,ENCRYPTED
DEK-Info: DES-EDE3-CBC,E33B4DB15A5DA39E

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
-----END RSA PRIVATE KEY-----	
`

	PRIVATE_KEY_JWE = `
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCaLICvAjZTz6dZ
GHzifkT02Ogqv91d4CYzH9lDnPxZArqdaMOLjYB87bPmW1HgMsSJxQ29iXXt79h0
nZ1ieiyhiPZTOtpgWjIpTlC2GgCvKoU0VmtspHsNHBV+7c79jY8pQyKhnESNkMwm
D5o4ttZ0vo94JYLIcHPUQVGFmqYywpIhtsCa9sYL2X3SIMULJ6Cq+YRrIcYvkT5O
ERmf78LWAvuyzk3KtjfKvgbNK/PVMd4Qlt4JlhDVhw9D/BY08Noj4HQcmvRr0OHv
RazIUNnkWc/yljQXEGcjiogjEhffrb489GK2WinPIe3dUhzB2l01NFAtl5CCmIq4
oUxjbL/xAgMBAAECggEAEDrNFR9DftEmSb/FLcajFa9byVxHPmGhQ9J/eZmmCuy5
nmZv51nvA7e2L3K/jL1raSFgT+qPiousGqKY9cvstNiJLbvNT0VPcH+5CWJ1Xfs4
6IiMaHRsK7YgISuBlbl4L5zMoXykzs5GX19XGc9Nhh6lVb2FJfqIhviVT143TrJ5
AQCz2unFKanbRDbf4kM57HsmaUwhRneNp/clp4VP+WP18HqqXRrid3QwfiPjah/x
ZKn52TeG4fYkaObsZt6eH7IMKnYIQy6eiFr9x+o3sGASyIoe6Atnomc+9qHcEQVH
fUDcgktkeP8/dgxxbKp5EUWmn8xKUnuZ7lPnuxZw0QKBgQDJJ21pHwOr9EuB4GL/
cLzVRXsiFlXG8k4aViGDppCkCWbvIO7gISr8Jz/GEmm41uUulwN7XYVfqOXp5/JA
kM4+PKg7h4M59VilX1UdyFTh2VhTXJxhrZr9DmrdOl5UW4hiMArIaFkG9q5xDLqj
gl1bCgV/UY8UqW9DxeGZP/fkBQKBgQDENdqVSAOTFM7ZeCyRuQa14qb7o9n3Dbwa
rCvJhx7xubl/6yqc1WUaYNyxE6cMM39aqc0HylSYCO6JeNupkOBDIc4QQ/hFpHyN
MkM6QF1aGCaDbYhxqxbNy7GYXC9qNFp6RdJFeQPRArBmZCtKJ8fmYpgEftxu+wro
jRULFPp7/QKBgQCY09vMgkPH4VN82X5dlMnjP2bN/yosfOvaFpZf76z09C/AfsT+
hDSkXy2Uz2iDhsGZJCMBF4y9oRUNIaxsYZhQsMUgdVS+NCmle0iv2ASlkvwIWdR2
Ye/fU5Tdf/srHGACOX33xU/eeo0OVx12HRXQlUyX7t9GU3S1iSJdzLwKAQKBgGxk
x+2CTB879o3jOtQCIHfoz5Di0u4N02X0yXfawriNfrHxS6g6p1DsQ987WSR/apK+
jXsJRrR68rRVZRyG4a2Uhk8sDYMDvJ8QLl2G40t7XgNrRl5tQvrL1b+y5arJY1Z7
Lg+dLAOSdbsLCXM9CMz4mLybDNHus/cGwaJOo5ZNAoGBAIVxS3sgC6HO06sKjyah
GnNp6Aa0ycG7mz7sprvXUhR3cwE8vUBwuDou3cDI9ojM/j2DxTp6XrK6Mpbyx686
WggiLmmjHyzYtNl8B2kFFKDY5f5yEn6jGxIbq3/2t5H9FmcP5EV7lh31qrJYMSc7
s/pi8RDMGUU5KCBEOJBvidP4
-----END PRIVATE KEY-----
`

	PUBLIC_KEY = `
-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAwv8Z+2+ULPFy7IXF1MeV
Z9wssTRNInF5zJ5dh+wpgyTG3IWZ2x7Y/XFVpvkJ9GndJVuKWjFvrS2hSeYovbe2
xPu6dvJl2P4O2UeJg27x+JkSySR74jdYGknqDq1VNPF1x7sV+T+g7cT9hyIO8QSs
VDW1uswnYDgfZpJw3rZtIyCF4OWZytTQ31thVYpGaWXoDWz6vSmLro8YtrBfIeZc
OzyHQM5an/mnt4TSYsHqX3hittfOyGckfBhb+zsySE3ZhjR9dWwvvgwq3sRjWnpk
p7JbWwdMnMcTxaQw8EPmo8YrF6ryVCeRXVduKBYY/yLbo9rFHBuLh1uV19pBzz8I
1VrLZF088O0iy9u+Dnjvf+2Gnhnv9Wc74rOm+rOiXG01yr0JtXgyQTvWK3jN6j7x
a5ZY+ZfFeOjgzv/o6M4JImQlmyHvx+6xcnk693KMfztCv/7gOaZW/UauB+FawIJ2
c+zZoYfAK8enr8/mJHFFZ0ohIdKmW1+HBeQgJsravaq6Awy37hucYqFM/NKgWW2o
8/V+kt3a1HLUMBf0v0sA9sW+6hoP22aJADXmxRWyisAyt2N1EEFTtdTcqkba+isV
yz6DZuviVu8/P8oLi+SvvTRLYaQIShd4QJOH7Uharq3FewWVg8ro+QMmdfdQUeHt
wLuWBqOmbXIe7QoNc3m8GKMCAwEAAQ==
-----END PUBLIC KEY-----
	
`

	TRUST_STORE = `
-----BEGIN CERTIFICATE-----
MIIFuzCCA6OgAwIBAgIUAuDeluPQ4r9XF/OUzk8G3goAVN0wDQYJKoZIhvcNAQEL
BQAwbTELMAkGA1UEBhMCTloxETAPBgNVBAgMCEF1Y2tsYW5kMREwDwYDVQQHDAhB
dWNrbGFuZDEQMA4GA1UECgwHOVNwb2tlczEMMAoGA1UECwwDRGV2MRgwFgYDVQQD
DA93d3cuOXNwb2tlcy5jb20wHhcNMjIwMTE5MjM1NDU4WhcNMzIwMTE3MjM1NDU4
WjBtMQswCQYDVQQGEwJOWjERMA8GA1UECAwIQXVja2xhbmQxETAPBgNVBAcMCEF1
Y2tsYW5kMRAwDgYDVQQKDAc5U3Bva2VzMQwwCgYDVQQLDANEZXYxGDAWBgNVBAMM
D3d3dy45c3Bva2VzLmNvbTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIB
AML/GftvlCzxcuyFxdTHlWfcLLE0TSJxecyeXYfsKYMkxtyFmdse2P1xVab5CfRp
3SVbiloxb60toUnmKL23tsT7unbyZdj+DtlHiYNu8fiZEskke+I3WBpJ6g6tVTTx
dce7Ffk/oO3E/YciDvEErFQ1tbrMJ2A4H2aScN62bSMgheDlmcrU0N9bYVWKRmll
6A1s+r0pi66PGLawXyHmXDs8h0DOWp/5p7eE0mLB6l94YrbXzshnJHwYW/s7MkhN
2YY0fXVsL74MKt7EY1p6ZKeyW1sHTJzHE8WkMPBD5qPGKxeq8lQnkV1XbigWGP8i
26PaxRwbi4dbldfaQc8/CNVay2RdPPDtIsvbvg5473/thp4Z7/VnO+Kzpvqzolxt
Ncq9CbV4MkE71it4zeo+8WuWWPmXxXjo4M7/6OjOCSJkJZsh78fusXJ5OvdyjH87
Qr/+4DmmVv1GrgfhWsCCdnPs2aGHwCvHp6/P5iRxRWdKISHSpltfhwXkICbK2r2q
ugMMt+4bnGKhTPzSoFltqPP1fpLd2tRy1DAX9L9LAPbFvuoaD9tmiQA15sUVsorA
MrdjdRBBU7XU3KpG2vorFcs+g2br4lbvPz/KC4vkr700S2GkCEoXeECTh+1IWq6t
xXsFlYPK6PkDJnX3UFHh7cC7lgajpm1yHu0KDXN5vBijAgMBAAGjUzBRMB0GA1Ud
DgQWBBQ08Y4X1jny1EnuGt7sZT/R0G2gkTAfBgNVHSMEGDAWgBQ08Y4X1jny1Enu
Gt7sZT/R0G2gkTAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4ICAQBg
CdCNIabpXf/TjSBWEK06oB4/ysUp+Wpzz6S+2wusc1rh1zXh2DqnU9VrhPUiDFec
D7eueh8Opr4VUYnWBltIBC+T1m+5vDLW7tYgdwH+AaZs4lwJ/jpZRdR7uFDrwlwJ
8Ko4GOuKMl4R7ZapFK6nccVPC1hOgQumXn64TvY1GgiYFg3BwDrOMbwvacQ0XJVA
XbAoHyvKuca765sPUlaQUDKXzzQvE2n5xGO70FlEljLiWOZ1PKud7uFNYh2JT2Ye
s/mjNEkmVJlANIMRCjBjvfCLCmETQ/rjYRjE0PG6m4oL44HsMcv1qhVkU3DKlNMi
uFSvHteo3W0BAKCCHXaHM1YUHNYNdOmjH7wVJW/irMqb0xd3qwTxJG7p4qAJs3up
wxiSqe2Tj0xmuyzeKeDiu/QRG6cqQ6ydtjtX/rCUhKx/Ux/D+XZSgcO1rn/5n5ig
FQvmyyNm4MApcowk+0D71m2RPpAgPR21uHlNfbo+AqnMaGR70EhUVIPKjUzf9VEa
uN5kxn/u2DgPhG5EoBJkQk4YFjziu5/0rQhS9brAlKRpa43NhvBQuunzDBQdtxVU
OtVXHzT1h0vTJUNXv1xIUfAx9tqqKMtPmyyaZhQkep7jWiec2ZCfxJ6eRmEo46eI
1nnwchBDYiUIvJdI3VssdEQXYk7Io8ZTDSuKdDQp0A==
-----END CERTIFICATE-----

-----BEGIN CERTIFICATE-----
MIIHTTCCBjWgAwIBAgIQS+ZmbZxdloWwrn/X34NGwjANBgkqhkiG9w0BAQsFADCB
ujELMAkGA1UEBhMCVVMxFjAUBgNVBAoTDUVudHJ1c3QsIEluYy4xKDAmBgNVBAsT
H1NlZSB3d3cuZW50cnVzdC5uZXQvbGVnYWwtdGVybXMxOTA3BgNVBAsTMChjKSAy
MDE0IEVudHJ1c3QsIEluYy4gLSBmb3IgYXV0aG9yaXplZCB1c2Ugb25seTEuMCwG
A1UEAxMlRW50cnVzdCBDZXJ0aWZpY2F0aW9uIEF1dGhvcml0eSAtIEwxTTAeFw0y
MTAzMDIyMTA5MTNaFw0yMjAzMDIyMTA5MTJaMIHmMQswCQYDVQQGEwJVUzERMA8G
A1UECBMISWxsaW5vaXMxEDAOBgNVBAcTB0NoaWNhZ28xEzARBgsrBgEEAYI3PAIB
AxMCVVMxGTAXBgsrBgEEAYI3PAIBAhMIRGVsYXdhcmUxJDAiBgNVBAoTG0Jhbmsg
b2YgQW1lcmljYSBDb3Jwb3JhdGlvbjEdMBsGA1UEDxMUUHJpdmF0ZSBPcmdhbml6
YXRpb24xEDAOBgNVBAUTBzI5Mjc0NDIxKzApBgNVBAMTIkF1dGhodWI0NjI1Mi1M
TEUuYmFua29mYW1lcmljYS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK
AoIBAQDCSfbC5jvkQ0LSov1GD3W11gRzq8q15k3264cTY/VsQbSOqhJdRCVVV8KV
QV2NhQW5eyij91k/cPfjcv38LqtEj7kxNVN7uRbluriKBfqm4d+3kUPkQ80IGvXy
cplmJtnlOKeugF3B11zmG2crdjzSr+gwNUP34BnZUELycotqK6Nl4edoWHRRG53f
VOlz+japwwhm5LM86Xa90yP1fFzxelNGbWNbS35S537D50ZRxq+P4blHhVTA+PXc
W4Jk1bcdQTNpAl9apeM2zp0OfBxj9ifVrknejmp+Nlo7IUcgwIWtPoFQx7uy1uTu
6RZIX1cpkUjeDztPzpHd3RH70I9rAgMBAAGjggMfMIIDGzAMBgNVHRMBAf8EAjAA
MB0GA1UdDgQWBBQPMs5yNndP7DGHMXeG8aO5pJAgODAfBgNVHSMEGDAWgBTD99C1
KjCtrw2RIXA5VN28iXDHOjBoBggrBgEFBQcBAQRcMFowIwYIKwYBBQUHMAGGF2h0
dHA6Ly9vY3NwLmVudHJ1c3QubmV0MDMGCCsGAQUFBzAChidodHRwOi8vYWlhLmVu
dHJ1c3QubmV0L2wxbS1jaGFpbjI1Ni5jZXIwMwYDVR0fBCwwKjAooCagJIYiaHR0
cDovL2NybC5lbnRydXN0Lm5ldC9sZXZlbDFtLmNybDAtBgNVHREEJjAkgiJBdXRo
aHViNDYyNTItTExFLmJhbmtvZmFtZXJpY2EuY29tMA4GA1UdDwEB/wQEAwIFoDAd
BgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwSwYDVR0gBEQwQjA3BgpghkgB
hvpsCgECMCkwJwYIKwYBBQUHAgEWG2h0dHBzOi8vd3d3LmVudHJ1c3QubmV0L3Jw
YTAHBgVngQwBATCCAX8GCisGAQQB1nkCBAIEggFvBIIBawFpAHcAVhQGmi/XwuzT
9eG9RLI+x0Z2ubyZEVzA75SYVdaJ0N0AAAF39MYKnAAABAMASDBGAiEA2XFl5C2a
/PUzRMJx3KSrZKyWZo9+FPkGSy49badxAsgCIQDCRKutrz4J7ACUHyMspC5NYFBB
rZB3hA2dT0ScBFAJzAB1AFWB1MIWkDYBSuoLm1c8U/DA5Dh4cCUIFy+jqh0HE9MM
AAABd/TGCrIAAAQDAEYwRAIgCz79wGXHCakpqSW9TdgKtZRVXLxAblWqixKq3UP9
tJ8CIHnV9CF5DZkMHvGN6PTZe2ZBtT5wNrsetzf2UfYrvRYVAHcARqVV63X6kSAw
taKJafTzfREsQXS+/Um4havy/HD+bUcAAAF39MYKvQAABAMASDBGAiEAkr+WV1fN
58d1W2zvQukTubZZRdq4NqYOgRW20krag6ACIQC6M31MwjsG1wX/feOP2SKWy2Ne
yOR2B1YjTXP52hwhCjANBgkqhkiG9w0BAQsFAAOCAQEAEISRbzs1mKd8Q7azDf//
ByDrZlqXUEn04cpUytiOAAcMbbSC7jxm0u5o+Zdz2ZGyDdzRa+FUwOdyPHoQHW/x
Hu36bRTvmgTvfvRVaJ9DL+BlcPJfEzhS+bFuuFKyngwibdU+I7GLouLY31AQZjQL
UGQ0iNKMCz501qbhk/JcJKXuzGdIepW4W+HROnLIJvZP6+Kmo5rhoVf2nscaUaVr
RdKTUcpeqPe+IjK+OX3F+HYUklSVsMhuwzI+c0PAvzEmEMkvb1MccCLiIR9Uxpiz
YtqJK4gkv/RZvMpzrlI7qRsFObi6hFck6aw3dj/7I2tAIOE31FjXtw0zkYdkHRac
YQ==
-----END CERTIFICATE-----
`
	CERT = `
MIIFuzCCA6OgAwIBAgIUAuDeluPQ4r9XF/OUzk8G3goAVN0wDQYJKoZIhvcNAQEL
BQAwbTELMAkGA1UEBhMCTloxETAPBgNVBAgMCEF1Y2tsYW5kMREwDwYDVQQHDAhB
dWNrbGFuZDEQMA4GA1UECgwHOVNwb2tlczEMMAoGA1UECwwDRGV2MRgwFgYDVQQD
DA93d3cuOXNwb2tlcy5jb20wHhcNMjIwMTE5MjM1NDU4WhcNMzIwMTE3MjM1NDU4
WjBtMQswCQYDVQQGEwJOWjERMA8GA1UECAwIQXVja2xhbmQxETAPBgNVBAcMCEF1
Y2tsYW5kMRAwDgYDVQQKDAc5U3Bva2VzMQwwCgYDVQQLDANEZXYxGDAWBgNVBAMM
D3d3dy45c3Bva2VzLmNvbTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIB
AML/GftvlCzxcuyFxdTHlWfcLLE0TSJxecyeXYfsKYMkxtyFmdse2P1xVab5CfRp
3SVbiloxb60toUnmKL23tsT7unbyZdj+DtlHiYNu8fiZEskke+I3WBpJ6g6tVTTx
dce7Ffk/oO3E/YciDvEErFQ1tbrMJ2A4H2aScN62bSMgheDlmcrU0N9bYVWKRmll
6A1s+r0pi66PGLawXyHmXDs8h0DOWp/5p7eE0mLB6l94YrbXzshnJHwYW/s7MkhN
2YY0fXVsL74MKt7EY1p6ZKeyW1sHTJzHE8WkMPBD5qPGKxeq8lQnkV1XbigWGP8i
26PaxRwbi4dbldfaQc8/CNVay2RdPPDtIsvbvg5473/thp4Z7/VnO+Kzpvqzolxt
Ncq9CbV4MkE71it4zeo+8WuWWPmXxXjo4M7/6OjOCSJkJZsh78fusXJ5OvdyjH87
Qr/+4DmmVv1GrgfhWsCCdnPs2aGHwCvHp6/P5iRxRWdKISHSpltfhwXkICbK2r2q
ugMMt+4bnGKhTPzSoFltqPP1fpLd2tRy1DAX9L9LAPbFvuoaD9tmiQA15sUVsorA
MrdjdRBBU7XU3KpG2vorFcs+g2br4lbvPz/KC4vkr700S2GkCEoXeECTh+1IWq6t
xXsFlYPK6PkDJnX3UFHh7cC7lgajpm1yHu0KDXN5vBijAgMBAAGjUzBRMB0GA1Ud
DgQWBBQ08Y4X1jny1EnuGt7sZT/R0G2gkTAfBgNVHSMEGDAWgBQ08Y4X1jny1Enu
Gt7sZT/R0G2gkTAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4ICAQBg
CdCNIabpXf/TjSBWEK06oB4/ysUp+Wpzz6S+2wusc1rh1zXh2DqnU9VrhPUiDFec
D7eueh8Opr4VUYnWBltIBC+T1m+5vDLW7tYgdwH+AaZs4lwJ/jpZRdR7uFDrwlwJ
8Ko4GOuKMl4R7ZapFK6nccVPC1hOgQumXn64TvY1GgiYFg3BwDrOMbwvacQ0XJVA
XbAoHyvKuca765sPUlaQUDKXzzQvE2n5xGO70FlEljLiWOZ1PKud7uFNYh2JT2Ye
s/mjNEkmVJlANIMRCjBjvfCLCmETQ/rjYRjE0PG6m4oL44HsMcv1qhVkU3DKlNMi
uFSvHteo3W0BAKCCHXaHM1YUHNYNdOmjH7wVJW/irMqb0xd3qwTxJG7p4qAJs3up
wxiSqe2Tj0xmuyzeKeDiu/QRG6cqQ6ydtjtX/rCUhKx/Ux/D+XZSgcO1rn/5n5ig
FQvmyyNm4MApcowk+0D71m2RPpAgPR21uHlNfbo+AqnMaGR70EhUVIPKjUzf9VEa
uN5kxn/u2DgPhG5EoBJkQk4YFjziu5/0rQhS9brAlKRpa43NhvBQuunzDBQdtxVU
OtVXHzT1h0vTJUNXv1xIUfAx9tqqKMtPmyyaZhQkep7jWiec2ZCfxJ6eRmEo46eI
1nnwchBDYiUIvJdI3VssdEQXYk7Io8ZTDSuKdDQp0A==
`
)

const DEX_KEYS = `{
	"keys": [
		{
			"use": "sig",
			"kty": "RSA",
			"kid": "d3467ed1c0518f3fd1b9e44e59354a902b878515",
			"alg": "RS256",
			"n": "ynoIb7pcIp9aV0cP-mmJfYj6OXHy3uULOUHtu14LsRLpSG5S2itgl1TWeCeF87ICc1wTPKl2yoti27iYWKX1tk6sUHpCDW2Fku5FP1SApQRiS7jGbN_QwXC9aFQlwaYk42FzzUar8EkPxn5cwdpe95qIq8Qz7Uhdgx47zMokZb2IBoTVIA1aPF3fUVhKrsa9UU2IYNfKb50-C87tQxkjcLrbMYgVveOB3xWUw8Cd9gwmPlHP8KTRx1UE4BYAMkm2HF9ZeA2NHhcJjc8ZM6lFMyQifyoAbGvjb5zE0i_HXnCekF4L35w51G57oeYIG4wFqeBmcmLe48prOzNZk22qkQ",
			"e": "AQAB"
		}
	]
}`

var keysTestServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(DEX_KEYS))
}))

func newTestContext(mode string) (*Context, error) {
	trustStore, err := os.CreateTemp(".", "*.crt")
	if err != nil {
		return nil, fmt.Errorf("failed to load the test trust store: %s", err.Error())
	}
	defer os.Remove(trustStore.Name())
	trustStore.WriteString(TRUST_STORE)
	trustStore.Close()

	privateKey, err := os.CreateTemp(".", "*.key")
	if err != nil {
		return nil, fmt.Errorf("failed to load the test private key: %s", err.Error())
	}
	defer os.Remove(privateKey.Name())
	secret := "secret"
	if mode == "jwe" {
		privateKey.WriteString(PRIVATE_KEY_JWE)
		secret = ""
	}
	privateKey.WriteString(PRIVATE_KEY_JWS)
	privateKey.Close()

	return New(keysTestServer.URL, "./"+trustStore.Name(), "./"+privateKey.Name(), secret)
}

func Test_ValidateJWS(t *testing.T) {

	ctx, err := newTestContext("jws")
	if err != nil {
		t.Fatalf("failed to create the test context: %s", err.Error())
	}

	tests := []struct {
		name  string
		ctx   *Context
		token string
		want  string
		err   string
	}{
		{
			name:  "Certificate expired",
			ctx:   ctx,
			token: JWS_X5C_A,
			want:  "",
			err:   "certificate is not trusted",
		},
		{
			name:  "Certificate not trusted",
			ctx:   &Context{},
			token: JWS_X5C_C,
			want:  "",
			err:   "certificate is not trusted",
		},
		{
			name:  "Token expired",
			ctx:   ctx,
			token: JWS_X5C_B,
			want:  "",
			err:   "Token is expired",
		},
		{
			name:  "Valid token",
			ctx:   ctx,
			token: JWS_X5C_C,
			want:  "123456789",
			err:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.ctx.Validate(tt.token)
			if err != nil && (tt.err == "" || !regexp.MustCompile(tt.err).MatchString(err.Error())) {
				t.Fatalf("unexpected error: got [%s], expecting [%s]", err.Error(), tt.err)
			}

			if err == nil && tt.err != "" {
				t.Fatalf("expecting error [%s], got none", tt.err)
			}

			if tt.err == "" && got["sub"] != tt.want {
				t.Fatalf("expecting subject [%s], got [%s]", tt.want, got["sub"])
			}

		})
	}
}

func Test_ValidateJWE(t *testing.T) {

	ctx, err := newTestContext("jwe")
	if err != nil {
		t.Fatalf("failed to create the test context: %s", err.Error())
	}

	tests := []struct {
		name  string
		ctx   *Context
		token string
		want  string
		err   string
	}{
		{
			name:  "Token is expired",
			ctx:   ctx,
			token: JWE,
			want:  "",
			err:   "Token is expired",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.ctx.Validate(tt.token)
			if err != nil && (tt.err == "" || !regexp.MustCompile(tt.err).MatchString(err.Error())) {
				t.Fatalf("unexpected error: got [%s], expecting [%s]", err.Error(), tt.err)
			}

			if err == nil && tt.err != "" {
				t.Fatalf("expecting error [%s], got none", tt.err)
			}

			if tt.err == "" && got["refresh_token"] != tt.want {
				t.Fatalf("expecting subject [%s], got [%s]", tt.want, got["sub"])
			}

		})
	}
}

func Test_fetchJWKS(t *testing.T) {

	tests := []struct {
		name    string
		jwksURL string
		err     bool
	}{
		{
			name:    "Successful fetch",
			jwksURL: keysTestServer.URL,
			err:     false,
		},
		{
			name:    "Failed fetch",
			jwksURL: "https://nsp-dev.9spokes.io/dex/keyz",
			err:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetchJWKS(tt.jwksURL)
			if err != nil && tt.err {
				return
			}

			if err == nil && tt.err {
				t.Fatalf("expecting error, got none")
			}
			if !tt.err {
				if err != nil {
					t.Fatalf("unexpected err: %s", err.Error())
				}
				if len(got) == 0 {
					t.Fatalf("expecting non empty key map")
				}
			}
		})
	}
}

func Test_getSigningKey(t *testing.T) {

	ctx, err := newTestContext("jws")
	if err != nil {
		t.Fatalf("failed to create the test context: %s", err.Error())
	}

	tests := []struct {
		name  string
		token string
		err   string
	}{
		{
			name:  "x5c not trusted",
			token: JWS_X5C_A,
			err:   "certificate is not trusted",
		},
		{
			name:  "x5c success",
			token: JWS_X5C_B,
			err:   "",
		},
		{
			name:  "invalid kid",
			token: JWS_KID,
			err:   "no key found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			token, _, err := new(jwt.Parser).ParseUnverified(tt.token, jwt.MapClaims{})
			if err != nil {
				t.Fatalf("failed to parse token: %s", err.Error())
			}

			key, err := ctx.getSigningKey(token)
			if err != nil && (tt.err == "" || !regexp.MustCompile(tt.err).MatchString(err.Error())) {
				t.Fatalf("unexpected error: got [%s], expecting [%s]", err.Error(), tt.err)
			}

			if err == nil && tt.err != "" {
				t.Fatalf("expecting error [%s], got none", tt.err)
			}

			if tt.err == "" && key == nil {
				t.Fatalf("expecting signing key, got none")
			}
		})
	}
}
