variable "center_name" {
  description = "Name of LeoStream Center"
  type        = string
  sensitive   = false
  default     = "aws-center-us-east-1"
}

variable "pool_image_name" {
  description = "Name of the image"
  type        = string
  sensitive   = false
  default     = "emr 5.23.0-ami-roller-7 hvm ebs"
}

variable "desktop_type" {
  description = "Type of desktop"
  type        = string
  sensitive   = false
  default     = "t2.micro"
}

variable "aws_sec_group_id" {
  description = "AWS security group"
  type        = string
  sensitive   = false
  default     = "sg-b1c127f6"
}

variable "aws_vpc_id" {
  description = "AWS security group"
  type        = string
  sensitive   = false
  default     = "vpc-9c880fe6"
}

variable "aws_iam_name" {
  description = "AWS IAM name"
  type        = string
  sensitive   = false
  default     = "leostream-desktop"
}

variable "aws_sub_net" {
  description = "AWS subnet"
  type        = string
  sensitive   = false
  default     = "subnet-09bd8643/vpc-9c880fe6/us-east-1d/,subnet-15daf71a/vpc-9c880fe6/us-east-1f/"
}

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
