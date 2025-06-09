locals {
  gsi_indexes = var.secondary_index_list != null ? [
    for idx in var.secondary_index_list : idx
    if lookup(idx, "type", "GSI") == "GSI"
  ] : []

  lsi_indexes = var.secondary_index_list != null ? [
    for idx in var.secondary_index_list : idx
    if lookup(idx, "type", "GSI") == "LSI"
  ] : []
}

resource "aws_dynamodb_table" "this" {
  name             = var.table_name
  billing_mode     = var.billing_mode
  hash_key         = var.hash_key
  range_key        = var.range_key
  stream_enabled   = var.stream_enabled
  stream_view_type = var.stream_type

  dynamic "attribute" {
    for_each = var.attributes
    content {
      name = attribute.value.name
      type = attribute.value.type
    }
  }

  dynamic "global_secondary_index" {
    for_each = local.gsi_indexes
    content {
      name               = global_secondary_index.value.name
      hash_key           = global_secondary_index.value.hash_key
      range_key          = lookup(global_secondary_index.value, "range_key", null)
      projection_type    = global_secondary_index.value.projection_type
      read_capacity      = global_secondary_index.value.read_capacity
      write_capacity     = global_secondary_index.value.write_capacity
      non_key_attributes = global_secondary_index.value.projection_type == "INCLUDE" ? global_secondary_index.value.non_key_attributes : null
    }
  }

  dynamic "local_secondary_index" {
    for_each = local.lsi_indexes
    content {
      name               = local_secondary_index.value.name
      range_key          = local_secondary_index.value.range_key
      projection_type    = local_secondary_index.value.projection_type
      non_key_attributes = local_secondary_index.value.projection_type == "INCLUDE" ? local_secondary_index.value.non_key_attributes : null
    }
  }

  dynamic "ttl" {
    for_each = var.ttl_enabled ? [1] : []
    content {
      attribute_name = var.ttl_attribute_name
      enabled        = true
    }
  }

  tags = merge(
    var.shared_tags,
    {
      "Name" = var.table_name,
    }
  )

  point_in_time_recovery {
    enabled = true
  }

  server_side_encryption {
    enabled = true
  }
}