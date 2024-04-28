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
tags: <default_tags for the provider>
  - name: <The tag key>
    value: <The tag value>
    ...


````


## Example for a project structure
````
ubuntu@hostname:~/project$ ../terraform_tree/tr -n 2
Enter a name: network
Enter a name: compute
2024/04/18 20:58:08 /home/ubuntu/project/modules/compute
2024/04/18 20:58:08 /home/ubuntu/project/modules/network

project/
├── main.tf
└── modules
    ├── compute
    │   ├── main.tf
    │   ├── outputs.tf
    │   └── variables.tf
    └── network
        ├── main.tf
        ├── outputs.tf
        └── variables.tf

````

## project/main.tf
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
  region = "ap-south-1"
}

module "network" {
  source = "./modules/network"
}

module "compute" {
  source = "./modules/compute"
}
````

