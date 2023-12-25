-- signup.lua
wrk.method = "POST"
wrk.body   = '{"email": "test@example.com", "password": "Password123!", "confirmPassword": "Password123!"}'
wrk.headers["Content-Type"] = "application/json"
