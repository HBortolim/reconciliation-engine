terraform {
  backend "s3" {
    bucket         = "reconciliation-engine-terraform-state"
    key            = "terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "reconciliation-engine-terraform-locks"
  }
}
