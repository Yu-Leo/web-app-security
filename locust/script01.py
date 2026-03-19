from locust import HttpUser, task, between
import random

class WebsiteUser(HttpUser):
    wait_time = between(1, 3)

    @task(3)
    def normal_request(self):
        self.client.get("/")

    @task(1)
    def suspicious_request(self):
        self.client.get("/admin?user=" + str(random.randint(1,1000)))
