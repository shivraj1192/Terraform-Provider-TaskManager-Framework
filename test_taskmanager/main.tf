terraform {
  required_providers {
    taskmanager = {
      source  = "local/taskmanager/taskmanager"
      version = "0.1.0"
    }
  }
}

provider "taskmanager" {
  base_url = "http://localhost:8080"
  token    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3NzkxODEyODgsInVzZXJfaWQiOjE0fQ.WJrTc9gZsRcSYFMB__zPfQR7S1Gb54dVhpO6OV3H5N4"
}

resource "taskmanager_user" "user_new" {
  uname    = "AT - TASKMANAGER UNAME"
  name     = "AT - TASKMANAGER NAME"
  email    = "AT.TASKMANAGER@gmail.com"
  password = "AT - TASKMANAGER PASSWORD"
  role     = "Member"
}

data "taskmanager_user" "data_user" {
    id = 12
}