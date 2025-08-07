variable "local_db_url" {
  type = string
}

variable "prod_db_url" {
  type = string

}

env "local" {
  src = "file://schema.hcl"
  dev = "docker://mysql/8/dev"

  url = var.local_db_url
  migration {
    dir = "file://migrations"
  }
}


env "prod" {
  src = "file://schema.hcl"
  dev = "docker://mysql/8/prod"

  url = var.prod_db_url
  migration {
    dir = "file://migrations"
  }
}