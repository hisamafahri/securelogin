package view

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hisamafahri/securelogin/internal/service"
)

type SigninView struct {
	requestID       string
	rayID           string
	ipAddress       string
	providers       []service.ProviderInfo
	applicationName string
}

func NewSigninView(
	requestID, rayID, ipAddress, applicationName string,
	providers []service.ProviderInfo,
) *SigninView {
	return &SigninView{
		requestID:       requestID,
		rayID:           rayID,
		ipAddress:       ipAddress,
		providers:       providers,
		applicationName: applicationName,
	}
}

func (v *SigninView) Render(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")

	var providerButtons strings.Builder
	if len(v.providers) == 0 {
		providerButtons.WriteString(`
            <div class="w-full bg-red-50 border border-red-300 p-4 rounded text-center">
               <p class="text-sm text-red-700 font-medium mb-1">No authentication provider available</p>
               <p class="text-xs text-red-600">Please contact your system administrator for assistance.</p>
            </div>
         `)
	} else {
		for _, provider := range v.providers {
			fmt.Fprintf(&providerButtons, `
            <form action="/signin/identifier" method="POST" enctype="multipart/form-data">
               <input type="hidden" name="provider_id" value="%s">
               <input type="hidden" name="request_id" value="%s">
               <button type="submit" class="w-full bg-white border border-gray-300 flex items-center justify-center gap-2 p-2 rounded text-gray-500 hover:bg-gray-50">
                  %s
                  %s
               </button>
            </form>
         `, provider.ID, v.requestID, provider.Icon, provider.Name)
		}
	}

	html := fmt.Sprintf(`
	<!DOCTYPE html>
<html lang="en">
   <head>
      <meta charset="UTF-8">
      <meta name="viewport" content="width=device-width, initial-scale=1.0">
      <title>Sign In | Secure Login</title>
      <script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"></script>
      <script src="https://cdn.tailwindcss.com"></script>
   </head>
   <body class="flex flex-col items-center justify-between h-screen bg-gray-100">
      <div></div>
      <div class="w-full w-full max-w-sm flex flex-1 flex-col items-center justify-center px-4">
         <div class="flex items-center justify-center gap-2 mb-6">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="size-8 text-purple-700" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75 11.25 15 15 9.75m-3-7.036A11.959 11.959 0 0 1 3.598 6 11.99 11.99 0 0 0 3 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285Z" />
            </svg>
            <p class="text-gray-700 text-2xl">Secure Login</p>
         </div>
         <div class="border-gray-300 border p-2 w-full h-fit bg-white rounded flex flex-col gap-2">
            <p class="text-sm text-gray-500 text-center px-4 mb-2">Enter your credentials to continue to:<br />
               <span class="font-medium text-gray-700">%s</span>
            </p>
            %s
         </div>
      </div>
      <div class="py-4 flex items-center justify-center">
         <p class="text-xs text-gray-500 text-center px-4">Ray ID: <span class="text-gray-700 font-medium">%s</span> • Your IP: <span class="text-gray-700 font-medium">%s</span></p>
      </div>
   </body>
</html>
`, v.applicationName, providerButtons.String(), v.rayID, v.ipAddress)

	c.String(http.StatusOK, html)
}
