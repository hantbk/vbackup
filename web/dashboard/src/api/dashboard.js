import request from '@/utils/request'

/**
 * Fetch homepage data
 * @param query
 * @returns {AxiosPromise}
 */
export function fetchIndex(query) {
  return request({
    url: '/dashboard/index',
    method: 'get',
    params: query
  })
}

/**
 * Fetch all repository statistics
 * @returns {AxiosPromise}
 */
export function fetchDoGetAllRepoStats() {
  return request({
    url: '/dashboard/doGetAllRepoStats',
    method: 'post'
  })
}

/**
 * Fetch operation logs
 * @param query
 * @returns {AxiosPromise}
 */
export function fetchLogs(query) {
  return request({
    url: '/dashboard/logs',
    method: 'get',
    params: query
  })
}
