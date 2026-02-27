export interface AuthCallbackRequest {
  session_state: string
  iss: string
  code: string
}