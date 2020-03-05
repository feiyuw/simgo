import React from 'react'
import axios from 'axios'
import { message, Button, Form, Upload, Icon, Input, InputNumber } from 'antd'
import {GRPC, FormItemLayoutWithOutLabel, TwoColumnsFormItemLayout} from '../constants'


class GrpcServerForm extends React.Component {
  protoFiles = []

  handleSubmit = e => {
    e.preventDefault()
    this.props.form.validateFields(async (err, values) => {
      if (err!==null && err!==undefined) {
        return
      }
      try{
        await axios.post('/api/v1/servers', {
          name: values.name,
          port: values.port,
          protocol: GRPC,
          options: {
            protos: values.protos.map(file => file.response.filepath),
          },
        })
      } catch (err) {
        return message.error('failed to add new server')
      }
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
          <Form.Item {...TwoColumnsFormItemLayout} label='Name'>
            {getFieldDecorator('name', {
              initialValue: '',
              rules: [{ required: true, message: 'Please input server name!' }],
            })(
              <Input style={{width: 350}} placeholder='type server name here' />
            )}
          </Form.Item>
          <Form.Item {...TwoColumnsFormItemLayout} label='Port'>
            {getFieldDecorator('port', {
              initialValue: 12345,
              rules: [{ required: true, message: 'Please input simulated server port!' }],
            })(
              <InputNumber style={{width: 350}} placeholder='type server port here' />
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
              Start
            </Button>
          </Form.Item>
        </Form>
  }
}


export default Form.create()(GrpcServerForm)
