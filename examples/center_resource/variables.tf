variable "center_password" {
  description = "AWS center password"
  type        = string
  sensitive   = true
}

variable "center_aws_key" {
  description = "AWS center key"
  type        = string
  sensitive   = true
}

variable "leostream_api_password" {
  description = "Leostream API password"
  type        = string
  sensitive   = true
}
