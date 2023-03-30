import service from '@/utils/request'

export const getTableApi = (data) => {
  return service({
    url: '/chatGpt/getTable',
    method: 'post',
    data
  })
}

export const createSKApi = (data) => {
  return service({
    url: '/chatGpt/createSK',
    method: 'post',
    data
  })
}

export const getSKApi = () => {
  return service({
    url: '/chatGpt/getSK',
    method: 'get',
  })
}

export const deleteSKApi = () => {
  return service({
    url: '/chatGpt/deleteSK',
    method: 'delete'
  })
}

export const testConnectApi = (data)=>{
  return service({
    url: '/chatGpt/testConnect',
    method: 'post',
    data
  })
}

export const getSchemaApi = (data)=>{
  return service({
    url: '/chatGpt/getSchema',
    method: 'post',
    data
  })
}
