import request from '@/utils/request'

/**
 * Searches for tasks based on the provided query parameters.
 * @param query - Query parameters for searching tasks
 * @returns {AxiosPromise}
 */
export function fetchSearch(query) {
  return request({
    url: '/task',
    method: 'get',
    params: query
  })
}

/**
 * Initiates an immediate backup based on the specified plan ID.
 * @param plan_id - ID of the backup plan to execute
 * @returns {AxiosPromise}
 */
export function fetchBackup(plan_id) {
  return request({
    url: `/task/backup/${plan_id}`,
    method: 'post'
  })
}

/**
 * Restores data from a specific snapshot within a repository.
 * @param repoid - ID of the repository to restore from
 * @param snapid - ID of the snapshot to restore
 * @param data - Additional data required for the restore operation
 * @returns {AxiosPromise}
 */
export function fetchRestore(repoid, snapid, data) {
  return request({
    url: `/task/${repoid}/restore/${snapid}`,
    method: 'post',
    data
  })
}
