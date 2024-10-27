import request from '@/utils/request'

/**
 * Fetches the list of repositories based on the provided query.
 * @param query - Query parameters for filtering repositories
 * @returns {AxiosPromise}
 */
export function fetchList(query) {
  return request({
    url: '/repository',
    method: 'get',
    params: query
  })
}

/**
 * Fetches details of a specific repository by its ID.
 * @param id - ID of the repository
 * @returns {AxiosPromise}
 */
export function fetchGet(id) {
  return request({
    url: `/repository/${id}`,
    method: 'get'
  })
}

/**
 * Deletes a specific repository by its ID.
 * @param id - ID of the repository
 * @returns {AxiosPromise}
 */
export function fetchDel(id) {
  return request({
    url: `/repository/${id}`,
    method: 'delete'
  })
}

/**
 * Creates a new repository with the provided data.
 * @param data - Data for the new repository
 * @returns {AxiosPromise}
 */
export function fetchCreate(data) {
  return request({
    url: '/repository',
    method: 'post',
    data
  })
}

/**
 * Updates an existing repository with the provided data.
 * @param data - Data to update the repository, including the repository ID
 * @returns {AxiosPromise}
 */
export function fetchUpdate(data) {
  return request({
    url: `/repository/${data.id}`,
    method: 'put',
    data
  })
}

/**
 * Fetches the list of snapshots for a specific repository.
 * @param query - Query parameters including the repository ID
 * @returns {AxiosPromise}
 */
export function fetchSnapshotsList(query) {
  return request({
    url: `/restic/${query.id}/snapshots`,
    method: 'get',
    params: query
  })
}

/**
 * Lists the contents of a specific snapshot in a repository.
 * @param repo - Repository ID
 * @param snap - Snapshot ID
 * @param query - Additional query parameters
 * @returns {AxiosPromise}
 */
export function fetchLsList(repo, snap, query) {
  return request({
    url: `/restic/${repo}/ls/${snap}`,
    method: 'get',
    params: query
  })
}

/**
 * Searches for specific data within a snapshot of a repository.
 * @param repo - Repository ID
 * @param snap - Snapshot ID
 * @param query - Additional query parameters
 * @returns {AxiosPromise}
 */
export function fetchSearchList(repo, snap, query) {
  return request({
    url: `/restic/${repo}/search/${snap}`,
    method: 'get',
    params: query
  })
}

/**
 * Fetches parameters for a specific repository.
 * @param id - Repository ID
 * @returns {AxiosPromise}
 */
export function fetchParmsList(id) {
  return request({
    url: `/restic/${id}/parms`,
    method: 'get'
  })
}

/**
 * Fetches parameters specific to the current user for a repository.
 * @param id - Repository ID
 * @returns {AxiosPromise}
 */
export function fetchParmsMyList(id) {
  return request({
    url: `/restic/${id}/parmsForMy`,
    method: 'get'
  })
}

/**
 * Checks the integrity of the specified repository.
 * @param repo - Repository ID
 * @returns {AxiosPromise}
 */
export function fetchCheck(repo) {
  return request({
    url: `/restic/${repo}/check`,
    method: 'post'
  })
}

/**
 * Rebuilds the index of the specified repository.
 * @param repo - Repository ID
 * @returns {AxiosPromise}
 */
export function fetchRebuildIndex(repo) {
  return request({
    url: `/restic/${repo}/rebuild-index`,
    method: 'post'
  })
}

/**
 * Prunes unused data from the specified repository to save space.
 * @param repo - Repository ID
 * @returns {AxiosPromise}
 */
export function fetchPrune(repo) {
  return request({
    url: `/restic/${repo}/prune`,
    method: 'post'
  })
}

/**
 * Upgrades the data format version of the specified repository.
 * @param repo - Repository ID
 * @returns {AxiosPromise}
 */
export function fetchMigrate(repo) {
  return request({
    url: `/restic/${repo}/migrate`,
    method: 'post'
  })
}

/**
 * Deletes a specific snapshot from the repository.
 * @param repo - Repository ID
 * @param snapshotid - ID of the snapshot to delete
 * @returns {AxiosPromise}
 */
export function fetchForget(repo, snapshotid) {
  return request({
    url: `/restic/${repo}/forget`,
    method: 'post',
    params: {
      snapshotid: snapshotid
    }
  })
}

/**
 * Unlocks the repository, optionally unlocking all locks.
 * @param repo - Repository ID
 * @param all - Boolean to unlock all locks
 * @returns {AxiosPromise}
 */
export function fetchUnlock(repo, all) {
  return request({
    url: `/restic/${repo}/unlock`,
    method: 'post',
    params: {
      all: all
    }
  })
}

/**
 * Fetches the last operation performed on the specified repository.
 * @param repo - Repository ID
 * @param type - Type of operation to fetch
 * @returns {AxiosPromise}
 */
export function fetchLastOper(repo, type) {
  return request({
    url: `/operation/last/${type}/${repo}`,
    method: 'get'
  })
}
