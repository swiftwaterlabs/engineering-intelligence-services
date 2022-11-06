resource "aws_glue_catalog_database" "object_db" {
  name = local.service_name
}

resource "aws_glue_catalog_table" "user" {
  name          = "repository"
  database_name = aws_glue_catalog_database.object_db.name

  table_type = "EXTERNAL_TABLE"

  parameters = {
    EXTERNAL              = "TRUE"
  }

  storage_descriptor {
    location      = "s3://${local.signal_bucket_name}/user"
    input_format  = "org.apache.hadoop.mapred.TextInputFormat"
    output_format = "org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat"

    ser_de_info {
      name                  = "json"
      serialization_library = "org.apache.hive.hcatalog.data.JsonSerDe"

      parameters = {
        "serialization.format" = 1
      }
    }

    columns {
      name = "id"
      type = "string"
    }

    columns {
      name = "host"
      type = "string"
    }

    columns {
      name = "organization"
      type = "string"
    }

    columns {
      name = "repository"
      type = "string"
    }

    columns {
      name = "Url"
      type = "string"
    }
  }
}