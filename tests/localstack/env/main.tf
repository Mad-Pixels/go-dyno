module "dynamodb_tables" {
  for_each = local.schemas

  source = "./dynamo"

  project              = "godyno-test"
  hash_key             = each.value.hash_key
  range_key            = each.value.range_key
  table_name           = each.value.table_name
  attributes           = each.value.attributes
  secondary_index_list = each.value.secondary_indexes

  shared_tags = local.tags
}