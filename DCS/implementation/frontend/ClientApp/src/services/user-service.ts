import http from '@/api/http'
import type { UserAllRequest, UserRolesByUserIdRequest } from '@/models/requests/user-request'
import type { UserAllResponse, UserRolesByUserIdResponse } from '@/models/responses/user-response'
import type { UserService as UserServiceI } from '@/models/services/user-service'
import type { UserProfile } from '@/models/user'
import type { UserRole } from '@/types/user-role'
import type { AxiosRequestHeaders, AxiosResponse } from 'axios'

const USER_BASE_URL = http.defaults.baseURL + '/users'

export const UserService: UserServiceI = {
  async getAllUsers(request?: UserAllRequest) {
    return Promise.resolve<AxiosResponse<UserAllResponse>>({
      data: { totalCount: mockUsers.length, items: mockUsers } as UserAllResponse,
      status: 200,
      statusText: 'OK',
      headers: {},
      config: {
        headers: {} as AxiosRequestHeaders,
      },
    }).then((res) => {
      console.log(res)
      return res.data.items
    })
    // return http
    //   .get<UserListResponse>('/users', { params: request, baseURL: USER_BASE_URL })
    //   .then((res) => res.data.items)
    //   .catch((err) => {
    //     console.error(err)
    //     return ({ totalCount: 0, items: [] } as UserListResponse).items
    //   })
  },

  async getRolesByUser(request: UserRolesByUserIdRequest) {
    return Promise.resolve<AxiosResponse<UserRolesByUserIdResponse>>({
      data: mockUsers.find((user) => user.id === request.userId)?.roleIds ?? [],
      status: 200,
      statusText: 'OK',
      headers: {},
      config: {
        headers: {} as AxiosRequestHeaders,
      },
    }).then((res) => {
      console.log(res)
      return res.data
    })
    // return http
    //   .get<UserRole[]>(`/users/${request.userId}/roles`, { baseURL: USER_BASE_URL })
    //   .then((res) => res.data)
    //   .catch((err) => {
    //     console.error(err)
    //     return [] as UserRole[]
    //   })
  },

  async getAuthorizedUsersWithRoles(...roles: [UserRole, ...UserRole[]]) {
    const allUsers = await this.getAllUsers()
    const authorizedUsers = await Promise.all(
      allUsers.map(async (user) => {
        const userRoles = await this.getRolesByUser({ userId: user.id })
        const isAuthorized = roles.some((role) => userRoles.includes(role))
        if (isAuthorized) {
          user.roleIds = roles.filter((role) => userRoles.includes(role))
        }
        return isAuthorized ? user : null
      }),
    )
    return authorizedUsers.filter((user) => user !== null)
  },
}

const mockUsers: UserProfile[] = [
  {
    participantId: 'part-001',
    firstName: 'John',
    lastName: 'Doe',
    email: 'john.doe@example.com',
    roleIds: ['TEMPLATE_APPROVER', 'TEMPLATE_CREATOR', 'TEMPLATE_MANAGER', 'TEMPLATE_REVIEWER'],
    id: '218979cc-47ac-4c14-88de-235bb18652a8',
    username: 'johndoe',
  },
  {
    participantId: 'part-002',
    firstName: 'Jane',
    lastName: 'Smith',
    email: 'jane.smith@example.com',
    roleIds: ['TEMPLATE_MANAGER', 'TEMPLATE_REVIEWER'],
    id: 'fcec22ac-e094-406b-930a-6ed3b7920414',
    username: 'janesmith',
  },
  {
    participantId: 'part-003',
    firstName: 'Bob',
    lastName: 'Johnson',
    email: 'bob.johnson@example.com',
    roleIds: ['TEMPLATE_APPROVER', 'TEMPLATE_MANAGER'],
    id: '4302ffb4-2dcb-4509-ae12-2863ec6b221b',
    username: 'bobjohnson',
  },
  {
    participantId: 'part-004',
    firstName: 'Alice',
    lastName: 'Williams',
    email: 'alice.williams@example.com',
    roleIds: ['TEMPLATE_REVIEWER'],
    id: '56775562-5fde-4396-ab37-0ae6d3812061',
    username: 'alicewilliams',
  },
  {
    participantId: 'part-005',
    firstName: 'Charlie',
    lastName: 'Brown',
    email: 'charlie.brown@example.com',
    id: '927bc415-f163-4282-bc89-614f0aa5da8f',
    username: 'charliebrown',
  },
  {
    participantId: 'part-006',
    firstName: 'Saoirse',
    lastName: 'Conrad',
    email: 'saoirse.conrad@example.com',
    roleIds: ['TEMPLATE_APPROVER', 'TEMPLATE_MANAGER'],
    id: '82d76203-6920-488a-9308-36984890bdd5',
    username: 'saoirseconrad',
  },
]
