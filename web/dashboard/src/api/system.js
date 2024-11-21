import request from '@/utils/request'

/**
 * Fetches a list of system files or directories based on the provided query parameters.
 * @param query - Query parameters for listing system files or directories
 * @returns {AxiosPromise}
 */
export function fetchLs(query) {
  return request({
    url: '/system/ls',
    method: 'get',
    params: query
  })
}

/**
 * Retrieves the current system version.
 * @returns {AxiosPromise}
 */
export function fetchVersion() {
  return request({
    url: '/system/version',
    method: 'get'
  })
}

/**
 * Fetches the latest available system version.
 * @returns {AxiosPromise}
 */
export function fetchLatestVersion() {
  return request({
    url: '/system/version/latest',
    method: 'get'
  })
}

/**
 * Upgrades the system to the specified version.
 * @param data - Version to upgrade to
 * @returns {AxiosPromise}
 */
export function fetchUpgradeVersion(data) {
  return request({
    url: '/system/upgradeVersion/' + data,
    method: 'post'
  })
}

/**
 * Retrieves the user's home directory path.
 * @returns {AxiosPromise}
 */
export function fetchUserHomePath() {
  return request({
    url: '/system/userHome',
    method: 'get'
  })
}
