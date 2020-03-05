const GRPC = 'grpc'
const HTTP = 'http'
const DUBBO = 'dubbo'
const PROTOCOLS= [
    {name: 'gRPC', value: 'grpc'},
    {name: 'HTTP', value: 'http'},
    {name: 'Dubbo', value: 'dubbo'},
  ]


const FormItemLayoutWithOutLabel = {
  wrapperCol: {
    xs: { span: 24, offset: 0 },
    sm: { span: 20, offset: 4 },
  },
}

const TwoColumnsFormItemLayout = {
  labelCol: {
    xs: { span: 20 },
    sm: { span: 8 },
  },
  wrapperCol: {
    xs: { span: 20 },
    sm: { span: 16 },
  },
}

export {
  PROTOCOLS,
  GRPC,
  HTTP,
  DUBBO,
  FormItemLayoutWithOutLabel,
  TwoColumnsFormItemLayout,
}
