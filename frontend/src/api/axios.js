import axios from 'axios'
import config from './config'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

export default function $axios(options) {
  return new Promise((resolve, reject) => {
    const instance = axios.create({
      // baseURL: config.baseUrl,
      headers: config.headers,
      timeout: config.timeout,
      withCredentials: config.withCredentials
    })

    instance.interceptors.request.use(
      config => {
        NProgress.start()
        return config
      },

      error => {
        console.log('request:', error)

        if (error.code === 'ECONNABORTED' && error.message.indexOf('timeout') !== -1) {
          console.log('timeout请求超时')
        }

        const errorInfo = error.response
        console.log("errorInfo", errorInfo)
        if (errorInfo) {
          error = errorInfo.data
          const errorStatus = errorInfo.status - line
          router.push({
            path: `/error/${errorStatus}`
          })
        }
        return Promise.reject(error)
      }
    )

    instance.interceptors.response.use(
      response => {
        NProgress.done()
        let data

        if (response.data === undefined) {
          data = JSON.parse(response.request.responseText)
        } else {
          data = response.data
        }

        return data
      },

      err => {
        NProgress.done()
        if (err && err.response) {
          switch (err.response.status) {
            case 400:
              err.response.message = '请求错误'
              break
            case 401:
              err.response.message = '未授权'
              break
            case 403:
              err.response.message = '拒绝访问'
              break
            case 404:
              err.response.message = `请求地址出错：${err.response.config.url}`
              break
            case 408:
              err.response.message = '请求超时'
              break
            case 500:
              err.response.message = '服务器内部错误'
              break
            case 501:
              err.response.message = '服务未实现'
              break
            case 502:
              err.response.message = '网关错误'
              break
            case 503:
              err.response.message = '服务不可用'
              break
            case 504:
              err.response.message = '网关超时'
              break
            case 505:
              err.response.message = 'HTTP版本不受支持'
              break
            default:
          }
        }
        return reject(err.response)
      }
    )

    instance(options).then(res => {
      resolve(res)
      return false
    }).catch(err => {
      reject(err)
    })
  })
}
