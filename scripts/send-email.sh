 curl -X POST http://localhost:8080/send-email \
     -H "Content-Type: application/json" \
     -H "X-Internal-API-Key: 1234567890" \
     -d '{
       "email": "dongw1023@gmail.com",
       "template": "verify_email",
       "data": {
         "name": "Local Tester",
         "verification_url": "https://localhost:3000/verify"
        }
      }'