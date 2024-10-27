import request from '@/utils/request'

/**
 * Fetches a list of policies based on the provided query parameters.
 * @param query - Query parameters for filtering policies
 * @returns {AxiosPromise}
 */
export function fetchList(query) {
  return request({
    url: '/policy',
    method: 'get',
    params: query
  })
}

/**
 * Creates a new policy with the provided data.
 * @param data - Data for the new policy
 * @returns {AxiosPromise}
 */
export function fetchCreate(data) {
  return request({
    url: '/policy',
    method: 'post',
    data
  })
}

/**
 * Updates an existing policy with the provided data.
 * @param data - Data to update the policy, including the policy ID
 * @returns {AxiosPromise}
 */
export function fetchUpdate(data) {
  return request({
    url: `/policy/${data.id}`,
    method: 'put',
    data
  })
}

/**
 * Deletes a specific policy by its ID.
 * @param id - ID of the policy
 * @returns {AxiosPromise}
 */
export function fetchDel(id) {
  return request({
    url: `/policy/${id}`,
    method: 'delete'
  })
}

/**
 * Executes a specific policy by its ID.
 * @param id - ID of the policy to execute
 * @returns {AxiosPromise}
 */
export function fetchDoPolicy(id) {
  return request({
    url: `/policy/do/${id}`,
    method: 'post'
  })
}
