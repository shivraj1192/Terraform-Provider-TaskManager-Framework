# Terraform-Provider-TaskManager-Framework

A custom Terraform provider for managing users in the TaskManager-Go backend using the Terraform Plugin Framework.

This provider is a new implementation of the TaskManager provider built with the Terraform Plugin Framework.  
Currently, it supports only:

- `taskmanager_user` resource
- `taskmanager_user` data source

New resources and data sources can be added later as the provider grows.

---

## Terraform Plugin Framework vs older SDK

Earlier Terraform providers were commonly built using Terraform Plugin SDKv2. The Terraform Plugin Framework is now the recommended approach for new providers because it gives a better development model and supports newer Terraform plugin protocol capabilities.

| Area                       | Terraform Plugin SDKv2            | Terraform Plugin Framework                              |
| -------------------------- | --------------------------------- | ------------------------------------------------------- |
| Development style          | Older SDK-based model             | Modern framework-based model                            |
| Type handling              | Less explicit in many areas       | Stronger typed values using framework types             |
| Diagnostics                | Error handling is less structured | Built-in diagnostics model                              |
| Provider structure         | More map/function oriented        | Interface-based provider/resource/data source structure |
| Terraform protocol support | Older provider model              | Recommended for protocol v5/v6 providers                |
| New provider development   | Still used in many providers      | Recommended for new providers                           |

The Plugin Framework helps developers focus more on API integration and provider behavior while the framework handles many Terraform-specific communication details.

---

## Table of Contents

- [Overview](#overview)
- [Current Features](#current-features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Step-by-Step Setup](#step-by-step-setup)
  - [1. Start the TaskManager-Go API](#1-start-the-taskmanager-go-api)
  - [2. Get Your API Token](#2-get-your-api-token)
  - [3. Clone and Build This Provider](#3-clone-and-build-this-provider)
  - [4. Install the Provider Locally](#4-install-the-provider-locally)
  - [5. Create Terraform Configuration](#5-create-terraform-configuration)
  - [6. Run Terraform Commands](#6-run-terraform-commands)
- [Provider Configuration](#provider-configuration)
- [Resource: taskmanager_user](#resource-taskmanager_user)
- [Data Source: taskmanager_user](#data-source-taskmanager_user)
- [Example main.tf](#example-maintf)
- [Testing](#testing)
- [Development](#development)
- [Acknowledgements](#acknowledgements)

---

## Overview

This Terraform provider allows you to manage users in the TaskManager-Go backend using Terraform.

The provider communicates with the TaskManager-Go API and performs user-related operations such as creating, reading, updating, and deleting users.

This version of the provider is built using the Terraform Plugin Framework, which is the newer recommended framework for developing Terraform providers.

> **Note:**  
> This provider requires the TaskManager-Go backend API to be running before using Terraform.

TaskManager backend repository:

```text
https://github.com/shivraj1192/TaskManager-Go
````

New provider repository:

```text
https://github.com/shivraj1192/Terraform-Provider-TaskManager-Framework/
```

---

## Current Features

Currently implemented features:

* Configure provider with TaskManager API base URL and JWT token
* Create a user using Terraform
* Read user details from the TaskManager backend
* Update user details
* Delete a user
* Fetch an existing user using a Terraform data source

---

## Project Structure

```text
terraform-provider-example/
│
├── main.go
├── taskmanager/
│	├── provider.go
│	├── datasources
│	│	└── datasource_user.go
│	├── resources
│	│	└── resource_user.go
│	└── helpers/
│		├── errs/
│		│	└── errors.go
│		└── imports/
│			└── imports.go
├── taskmanager_client
│   ├── client.go
│	├── plans.go
│   ├── models.go
│   ├── globals.go
│   └── user.go
├── go.mod
├── go.sum
└── test_taskmanager/
    └── main.tf
```

---

## Prerequisites

Before using this provider, make sure you have the following installed:

* Go 1.20 or later
* Terraform v1.0 or later
* Running TaskManager-Go backend API
* Valid JWT token from the TaskManager-Go API

---

## Step-by-Step Setup

### 1. Start the TaskManager-Go API

First, clone and run the TaskManager-Go backend application.

```sh
git clone https://github.com/shivraj1192/TaskManager-Go.git
cd TaskManager-Go

go mod tidy
go run ./cmd/main.go
```

By default, the backend should run on:

```text
http://localhost:8080/
```

Make sure the API is running before using this Terraform provider.

---

### 2. Get Your API Token

The provider uses JWT authentication.
You need to register or login through the TaskManager-Go API and get a token.

#### Register a user

Send a `POST` request to:

```text
http://localhost:8080/api/register
```

Example request body:

```json
{
  "uname": "adminuser",
  "name": "Admin User",
  "email": "admin@example.com",
  "password": "password123"
}
```

#### Login to get token

Send a `POST` request to:

```text
http://localhost:8080/api/login
```

Example request body:

```json
{
  "email": "admin@example.com",
  "password": "password123"
}
```

The response will contain a JWT token.

Copy that token and use it in your Terraform provider configuration.

---

### 3. Clone and Build This Provider

Clone the Terraform Plugin Framework provider repository:

```sh
git clone https://github.com/shivraj1192/Terraform-Provider-TaskManager-Framework.git
cd Terraform-Provider-TaskManager-Framework
```

Install dependencies:

```sh
go mod tidy
```

Build the provider binary.

#### Windows

```sh
go build -o terraform-provider-taskmanager.exe
```

#### Linux / Mac

```sh
go build -o terraform-provider-taskmanager
```

---

### 4. Install the Provider Locally

Because this provider is not published in the Terraform Registry, install it manually in the local Terraform plugin directory.

The local provider source used in Terraform configuration is:

```text
local/taskmanager/taskmanager
```

#### Windows

Create the plugin directory:

```powershell
$dest = "$env:APPDATA\terraform.d\plugins\local\taskmanager\taskmanager\0.1.0\windows_amd64"
New-Item -ItemType Directory -Force -Path $dest
```

Move the provider binary:

```powershell
Move-Item -Path ".\terraform-provider-taskmanager.exe" -Destination "$dest\terraform-provider-taskmanager.exe"
```

#### Linux

```sh
mkdir -p ~/.terraform.d/plugins/local/taskmanager/taskmanager/0.1.0/linux_amd64
mv terraform-provider-taskmanager ~/.terraform.d/plugins/local/taskmanager/taskmanager/0.1.0/linux_amd64/
```

#### Mac

For Intel Mac:

```sh
mkdir -p ~/.terraform.d/plugins/local/taskmanager/taskmanager/0.1.0/darwin_amd64
mv terraform-provider-taskmanager ~/.terraform.d/plugins/local/taskmanager/taskmanager/0.1.0/darwin_amd64/
```

For Apple Silicon Mac:

```sh
mkdir -p ~/.terraform.d/plugins/local/taskmanager/taskmanager/0.1.0/darwin_arm64
mv terraform-provider-taskmanager ~/.terraform.d/plugins/local/taskmanager/taskmanager/0.1.0/darwin_arm64/
```

#### To delete the plugin (cleanup):

Windows (PowerShell):
```powershell
Remove-Item -Path "$env:APPDATA\terraform.d\plugins\local\taskmanager\taskmanager\0.1.0\windows_amd64" -Recurse -Force
```
Linux:
```sh
rm -rf ~/.terraform.d/plugins/local/taskmanager/taskmanager/0.1.0/linux_amd64
```
Mac:
```sh
rm -rf ~/.terraform.d/plugins/local/taskmanager/taskmanager/0.1.0/darwin_amd64
```


---

### 5. Create Terraform Configuration

Create a new folder for testing the provider:

```sh
mkdir taskmanager-provider-test
cd taskmanager-provider-test
```

Create a `main.tf` file.

```hcl
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
```

---

### 6. Run Terraform Commands

Initialize Terraform:

```sh
terraform init
```

Check the execution plan:

```sh
terraform plan
```

Apply the configuration:

```sh
terraform apply
```

Destroy created resources:

```sh
terraform destroy
```

---

## Provider Configuration

The provider requires the following configuration:

```hcl
provider "taskmanager" {
  base_url = "http://localhost:8080/"
  token    = var.taskmanager_token
}
```

### Arguments

| Argument   | Required | Description                                |
| ---------- | -------- | ------------------------------------------ |
| `base_url` | Yes      | Base URL of the running TaskManager-Go API |
| `token`    | Yes      | JWT token used for API authentication      |

### Recommended Token Usage

Do not hardcode the token directly in `main.tf`.

Use an environment variable:

#### Windows Command Prompt

```cmd
set TF_VAR_taskmanager_token=PASTE_YOUR_TOKEN_HERE
```

#### Windows PowerShell

```powershell
$env:TF_VAR_taskmanager_token="PASTE_YOUR_TOKEN_HERE"
```

#### Linux / Mac

```sh
export TF_VAR_taskmanager_token=PASTE_YOUR_TOKEN_HERE
```

Or create a `terraform.tfvars` file:

```hcl
taskmanager_token = "PASTE_YOUR_TOKEN_HERE"
```

> Do not commit `terraform.tfvars` if it contains a real token.

---

## Resource: taskmanager_user

The `taskmanager_user` resource is used to create and manage users in the TaskManager-Go backend.

### Example

```hcl
resource "taskmanager_user" "example" {
  uname    = "john123"
  name     = "John Doe"
  email    = "john@example.com"
  password = "password123"
}
```

### Arguments

| Argument   | Required | Description                          |
| ---------- | -------- | ------------------------------------ |
| `uname`    | Yes      | Unique username of the user          |
| `name`     | Yes      | Full name of the user                |
| `email`    | Yes      | Email address of the user            |
| `password` | Yes      | Password used when creating the user |

### Attributes

| Attribute | Description                                |
| --------- | ------------------------------------------ |
| `id`      | User ID returned by the TaskManager-Go API |
| `uname`   | Username of the user                       |
| `name`    | Full name of the user                      |
| `email`   | Email address of the user                  |

### Import

If import support is implemented in the resource, an existing user can be imported using:

```sh
terraform import taskmanager_user.example USER_ID
```

---

## Data Source: taskmanager_user

The `taskmanager_user` data source is used to fetch an existing user from the TaskManager-Go backend.

### Example

```hcl
data "taskmanager_user" "example" {
  id = "USER_ID"
}
```

### Output Example

```hcl
output "user_name" {
  value = data.taskmanager_user.example.name
}

output "user_email" {
  value = data.taskmanager_user.example.email
}
```

### Arguments

| Argument | Required | Description             |
| -------- | -------- | ----------------------- |
| `id`     | Yes      | ID of the user to fetch |

### Attributes

| Attribute | Description               |
| --------- | ------------------------- |
| `id`      | User ID                   |
| `uname`   | Username of the user      |
| `name`    | Full name of the user     |
| `email`   | Email address of the user |

---

## Example main.tf

Below is a complete example using both the user resource and user data source.

```hcl
terraform {
  required_providers {
    taskmanager = {
      source  = "local/taskmanager/taskmanager"
      version = "0.1.0"
    }
  }
}

provider "taskmanager" {
  base_url = "http://localhost:8080/"
  token    = var.taskmanager_token
}

variable "taskmanager_token" {
  description = "JWT token for TaskManager-Go API"
  type        = string
  sensitive   = true
}

resource "taskmanager_user" "example" {
  uname    = "john123"
  name     = "John Doe"
  email    = "john@example.com"
  password = "password123"
}

data "taskmanager_user" "example" {
  id = taskmanager_user.example.id
}

output "created_user_id" {
  value = taskmanager_user.example.id
}

output "created_user_email" {
  value = data.taskmanager_user.example.email
}
```

Run:

```sh
terraform init
terraform plan
terraform apply
```

---

## Testing

Run all tests:

```sh
go test ./...
```

Run tests with verbose output:

```sh
go test -v ./...
```

If acceptance tests require the TaskManager-Go API, make sure the backend is running before executing tests.

---

## Development

This provider is built using the Terraform Plugin Framework.

Useful development commands:

```sh
go mod tidy
go fmt ./...
go test ./...
go build -o terraform-provider-taskmanager
```

Recommended development flow:

1. Start the TaskManager-Go backend.
2. Build the provider binary.
3. Install the binary in the local Terraform plugin directory.
4. Run `terraform init`.
5. Test resource and data source behavior using Terraform configuration.
6. Update provider code and rebuild when changes are made.

---

## Cleanup

To remove the locally installed provider binary:

### Windows

```powershell
Remove-Item -Path "$env:APPDATA\terraform.d\plugins\local\taskmanager\taskmanager\0.1.0" -Recurse -Force
```

### Linux

```sh
rm -rf ~/.terraform.d/plugins/local/taskmanager/taskmanager/0.1.0/linux_amd64
```

### Mac

```sh
rm -rf ~/.terraform.d/plugins/local/taskmanager/taskmanager/0.1.0/darwin_amd64
rm -rf ~/.terraform.d/plugins/local/taskmanager/taskmanager/0.1.0/darwin_arm64
```

---

## Acknowledgements

* Terraform
* Terraform Plugin Framework
* TaskManager-Go backend

---

This project is created for learning and understanding how to build a Terraform provider using the Terraform Plugin Framework.

````

Small correction you should make if your actual schema is different: if your data source fetches user by `email` instead of `id`, replace this part:

```hcl
data "taskmanager_user" "example" {
  id = taskmanager_user.example.id
}
````

with:

```hcl
data "taskmanager_user" "example" {
  email = taskmanager_user.example.email
}
```