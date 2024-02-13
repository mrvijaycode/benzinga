
## HOW TO RUN:

###### To start the application execute the below command
```./start.sh```

It will build and run the application.

You can test in postman

GET: ``localhost:8080/healthz``

![GET/healthz](<images/Screenshot 2024-02-14 at 12.54.12 AM.png>)

POST: ``localhost:8080/log``

![log](<images/Screenshot 2024-02-14 at 12.54.47 AM.png>)

The application will POST the logs to below webhook.site. Webhook details are given in .env file

![alt text](<images/Screenshot 2024-02-14 at 1.11.47 AM.png>)