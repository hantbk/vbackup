import request from '@/utils/request'

/**
 * Logs in a user with provided credentials.
 * @param data - User login credentials (e.g., username and password)
 * @returns {AxiosPromise} - The response from the server
 */
export function login(data) {
  return request({
    url: '/login',
    method: 'post',
    data
  })
}

/**
 * Refreshes the authentication token.
 * @returns {AxiosPromise} - The new authentication token
 */
export function refreshToken() {
  return request({
    url: '/refreshToken',
    method: 'post'
  })
}

/**
 * Retrieves information about the currently authenticated user.
 * @param token - The user's authentication token
 * @returns {Promise} - Mocked user data including roles and name
 */
export function getInfo(token) {
  return Promise.resolve({
    roles: ['admin'],
    introduction: 'I am a super administrator',
    name: 'Super Admin'
  })
}

/**
 * Fetches a list of all users.
 * @returns {AxiosPromise} - A list of user data
 */
export function fetchList() {
  return request({
    url: '/user',
    method: 'get'
  })
}

/**
 * Deletes a user by ID.
 * @param id - The ID of the user to be deleted
 * @returns {AxiosPromise} - The response from the server after deletion
 */
export function fetchDel(id) {
  return request({
    url: `/user/${id}`,
    method: 'delete'
  })
}

/**
 * Creates a new user with the provided data.
 * @param data - Data for the new user (e.g., name, email, role)
 * @returns {AxiosPromise} - The response from the server after creation
 */
export function fetchCreate(data) {
  return request({
    url: '/user',
    method: 'post',
    data
  })
}

/**
 * Updates an existing user's information.
 * @param data - Updated user data including user ID
 * @returns {AxiosPromise} - The response from the server after updating
 */
export function fetchUpdate(data) {
  return request({
    url: `/user/${data.id}`,
    method: 'put',
    data
  })
}

/**
 * Resets the password for a specific user.
 * @param data - Data required for password reset (e.g., user ID, new password)
 * @returns {AxiosPromise} - The response from the server after resetting the password
 */
export function fetchRePwd(data) {
  return request({
    url: '/repwd',
    method: 'post',
    data
  })
}

/**
 * Fetches the one-time password (OTP) configuration for the user.
 * @returns {AxiosPromise} - OTP configuration details
 */
export function fetchOtp() {
  return request({
    url: '/otp',
    method: 'get'
  })
}

/**
 * Binds a new OTP to the user's account.
 * @param data - OTP data required for binding (e.g., OTP code)
 * @returns {AxiosPromise} - The response from the server after binding OTP
 */
export function fetchBindOtp(data) {
  return request({
    url: '/otp',
    method: 'post',
    data
  })
}

/**
 * Deletes or unbinds the current OTP configuration from the user's account.
 * @returns {AxiosPromise} - The response from the server after deleting OTP
 */
export function fetchDeleteOtp() {
  return request({
    url: '/otp',
    method: 'put'
  })
}
