resource "aws_dynamodb_table" "hosts" {
  name           = "${local.service_name}_hosts"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "Id"

  attribute {
    name = "Id"
    type = "S"
  }

}

resource "aws_dynamodb_table" "event_sources" {
  name           = "${local.service_name}_event_sources"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "Id"

  attribute {
    name = "Id"
    type = "S"
  }

}