/**
 * Repository Types
 * @type {[{code: number, name: string}, {code: number, name: string}, null]}
 */
export const repoTypeList = [
  { code: 1, name: 'S3', tips: 'Minio' },
  { code: 3, name: 'Sftp', tips: 'Not recommended' },
  { code: 4, name: 'Local', tips: 'Recommended for testing' },
  { code: 5, name: 'Rest', tips: 'Recommended' }
]

/**
 * Retention Policy Types
 * @type {[{code: number, name: string, tips: string},{code: string, name: string, tips: string},{code: number, name: string, tips: string},{code: number, name: string, tips: string},{code: number, name: string, tips: string},null]}
 */
export const ForgetTypeList = [
  { code: 'last', name: 'Count', tips: 'Count' },
  { code: 'hourly', name: 'Hourly', tips: 'Hourly' },
  { code: 'daily', name: 'Daily', tips: 'Daily' },
  { code: 'weekly', name: 'Weekly', tips: 'Weekly' },
  { code: 'monthly', name: 'Monthly', tips: 'Monthly' },
  { code: 'yearly', name: 'Yearly', tips: 'Yearly' }
]

/**
 * Repository Connection Status
 * @type {[{code: number, color: string, name: string}, {code: number, color: string, name: string}, {code: number, color: string, name: string}]}
 */
export const repoStatusList = [
  { code: 1, name: 'Fetching', color: 'info' },
  { code: 2, name: 'Normal', color: 'success' },
  { code: 3, name: 'Error', color: 'danger' }
]

/**
 * Log Levels
 * @type {[{code: number, color: string, name: string}, {code: number, color: string, name: string}]}
 */
export const LoglevelList = [
  { code: 1, name: 'Info', color: 'info' },
  { code: 2, name: 'Warning', color: 'warning' },
  { code: 3, name: 'Success', color: 'success' },
  { code: 4, name: 'Error', color: 'error' }
]

/**
 * Compression Levels
 * @type {[{code: number, color: string, name: string},{code: number, color: string, name: string},{code: number, color: string, name: string}]}
 */
export const compressionList = [
  { code: 0, name: 'Auto', color: 'success' },
  { code: 1, name: 'Off', color: 'info' },
  { code: 2, name: 'Maximum', color: 'warning' }
]
