module "bite_tracker_eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "21.1.5"

  name                   = local.name
  kubernetes_version     = "1.33"
  region                 = local.region
  endpoint_public_access = true
  deletion_protection    = false

  enable_cluster_creator_admin_permissions = true

  # TODO: Allow access only to bite-tracker ECR

  compute_config = {
    enabled    = true
    node_pools = [local.name]
  }

  vpc_id     = aws_vpc.main.id
  subnet_ids = [aws_subnet.bite_tracker_k8s_a.id, aws_subnet.bite_tracker_k8s_b.id]

  tags = local.default_tags
}

# resource "aws_iam_role" "bite_tracker_cluster" {
#   name = "${local.name}-cluster"
#
#   assume_role_policy = jsonencode({
#     Version = "2012-10-17"
#     Statement = [
#       {
#         Action = ["sts:AssumeRole", "sts:TagSession"]
#         Effect = "Allow"
#         Principal = {
#           Service : [
#             "eks.amazonaws.com"
#           ]
#         }
#       },
#     ]
#   })
#   tags = local.default_tags
# }
#
# resource "aws_iam_policy" "bite_tracker_cluster_role" {
#   name = "${local.name}-cluster"
#   policy = jsonencode({
#     Version = "2012-10-17",
#     Statement = [
#       {
#         Effect = "Allow",
#         Action = [
#           "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
#           "arn:aws:iam::aws:policy/AmazonEKSComputePolicy",
#           "arn:aws:iam::aws:policy/AmazonEKSBlockStoragePolicy",
#           "arn:aws:iam::aws:policy/AmazonEKSLoadBalancingPolicy",
#           "arn:aws:iam::aws:policy/AmazonEKSNetworkingPolicy",
#           "arn:aws:iam::aws:policy/AmazonEKSServicePolicy",
#           "arn:aws:iam::aws:policy/AmazonEKSVPCResourceController",
#         ],
#         Resource = "*"
#       },
#     ]
#   })
#   tags = local.default_tags
# }
