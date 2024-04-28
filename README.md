# tfinit
The appliction creates a directory with tf files for each module. The following files will be created in the directories:
- main.tf
- variables.tf
- outputs.tf
And in the project "root" directory[^1] it creates:
- main.tf
- variables.tf
- terraform.tfvars

[^1]: The directory in which the appliction was executed

The program searches for a file called config.yml/yaml and retrieves the project's structure from it.
The YAML file format should be as follows:
````yaml
# required
region: <The aws region>(Example :us-east-1)

modules:
  <The module name>:
    vars:
      - name: <The variable name>
        type: <The variable type>
        description: <optional>


# optional
# default_tags for the provider
tags:
  - name: <The tag key>
    value: <The tag value>



````


## Example for a project structure
````
project/
├── config.yml
├── main.tf
├── modules
│   ├── compute
│   │   ├── main.tf
│   │   ├── outputs.tf
│   │   └── variables.tf
│   └── network
│       ├── main.tf
│       ├── outputs.tf
│       └── variables.tf
├── terraform.tfvars
└── variables.tf

````
### config.yaml
````yaml
region: us-east-1

tags:
  - name: Environment
    value: Test
  - name: project
    value: test

modules:
  network:
    vars:
      - name: availability_zones
        type: list(string)

  compute:
    vars:
      - name: instence_type
        type: string
        description: The instance type to use for the deployment

      - name: number_of_instances
        type: number
````

### project/main.tf
````hcl
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
  default_tags {
    tags = {
      Environment = "Test"
      project     = "test"
    }
  }
}

module "network" {
  source            = "./modules/network"
  availability_zones = var.availability_zones
}

module "compute" {
  source              = "./modules/compute"
  instence_type       = var.instence_type
  number_of_instances = var.number_of_instances
}
````

### project/variables.tf
````hcl
variable "availability_zones" {
  type = list(string)
}

variable "instence_type" {
  type        = string
  description = "The instance type to use for the deployment"
}
variable "number_of_instances" {
  type = number
}
````

### project/terraform.tfvars
````hcl
availability_zones = "placeholder"

instence_type       = "placeholder"
number_of_instances = "placeholder"
````

### modules/network/variables.tf
````hcl
variable "availability_zones" {
  type = list(string)
}
````
### modules/compute/variables.tf
````hcl
variable "instence_type" {
  type        = string
  description = "The instance type to use for the deployment"
}
variable "number_of_instances" {
  type = number
}
````
