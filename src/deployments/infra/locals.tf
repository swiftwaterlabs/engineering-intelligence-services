locals {
  service_name = "engineering_intelligence_${var.environment}"
  signal_bucket_name = replace("${local.service_name}data","_","")
  athena_bucket_name = replace("${local.service_name}athena","_","")
}