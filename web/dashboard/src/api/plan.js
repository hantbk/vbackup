import request from '@/utils/request'

/**
 * Creates a new plan with the provided data.
 * @param data - Data for the new plan
 * @returns {AxiosPromise}
 */
export function fetchCreate(data) {
  return request({
    url: '/plan',
    method: 'post',
    data
  })
}

/**
 * Updates an existing plan with the provided data.
 * @param data - Data to update the plan, including the plan ID
 * @returns {AxiosPromise}
 */
export function fetchUpdate(data) {
  return request({
    url: `/plan/${data.id}`,
    method: 'put',
    data
  })
}

/**
 * Deletes a specific plan by its ID.
 * @param id - ID of the plan
 * @returns {AxiosPromise}
 */
export function fetchDel(id) {
  return request({
    url: `/plan/${id}`,
    method: 'delete'
  })
}

/**
 * Fetches a list of plans based on the provided query parameters.
 * @param query - Query parameters for filtering plans
 * @returns {AxiosPromise}
 */
export function fetchList(query) {
  return request({
    url: '/plan',
    method: 'get',
    params: query
  })
}

/**
 * Fetches the next scheduled time for a plan based on the provided query parameters.
 * @param query - Query parameters to determine the next time
 * @returns {AxiosPromise}
 */
export function fetchNextTime(query) {
  return request({
    url: '/plan/next_time',
    method: 'get',
    params: query
  })
}
