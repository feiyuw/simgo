import React from 'react'
import axios from 'axios'
import { Button, Form, Upload, Icon, Input } from 'antd'


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

class GrpcClientForm extends React.Component {
  protoFiles = []

  handleSubmit = e => {
    e.preventDefault()
    this.props.form.validateFields(async (err, values) => {
      if (err!==null && err!==undefined) {
        return
      }
      await axios.post('/api/v1/clients', {
        server: values.server,
        protos: values.protos.map(file => file.response.filepath),
      })
      this.props.onSubmit()
    })
  }

  removeFile = async (file) => {
    await axios.delete(`/api/v1/files?filepath=${file.response.filepath}`)
  }

  normFile = e => {
    if (Array.isArray(e)) {
      return e
    }
    return e && e.fileList
  }

  render() {
    const { getFieldDecorator } = this.props.form

    return <Form onSubmit={this.handleSubmit}>
          <Form.Item {...TwoColumnsFormItemLayout} label='Server Address'>
            {getFieldDecorator('server', {
              initialValue: '',
              rules: [{ required: true, message: 'Please input server addr!' }],
            })(
              <Input style={{width: 350}} placeholder='type server addr here' />
            )}
          </Form.Item>
          <Form.Item {...TwoColumnsFormItemLayout} label='Proto Files'>
            {getFieldDecorator('protos', {
              valuePropName: 'fileList',
              getValueFromEvent: this.normFile,
              initialValue: [],
              rules: [{required: true, message: 'Please upload proto files!'}]
            })(
              <Upload
                accept='.proto'
                action='/api/v1/files'
                onRemove={this.removeFile}
                showUploadList={{showRemoveIcon: true, showDownloadIcon: false}}
              >
                <Button><Icon type='upload' /></Button>
              </Upload>
            )}
          </Form.Item>
          <Form.Item {...FormItemLayoutWithOutLabel}>
            <Button type="primary" htmlType="submit">
              Connect
            </Button>
          </Form.Item>
        </Form>
  }
}


export default Form.create()(GrpcClientForm)
