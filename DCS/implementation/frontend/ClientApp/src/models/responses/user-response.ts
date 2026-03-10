import type { UserRole } from "@/types/user-role"
import type { UserProfile } from "../user-profile"

export interface UserAllResponse {
  totalCount: number
  items: UserProfile[]
}

export type UserRolesByUserIdResponse = UserRole[]
