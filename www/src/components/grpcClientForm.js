import React from 'react'
import axios from 'axios'
import { message, Button, Form, Upload, Icon, Input } from 'antd'
import {GRPC, FormItemLayoutWithOutLabel, TwoColumnsFormItemLayout} from '../constants'
import urls from '../urls'


class GrpcClientForm extends React.Component {
  protoFiles = []

  handleSubmit = e => {
    e.preventDefault()
    this.props.form.validateFields(async (err, values) => {
      if (err!==null && err!==undefined) {
        return
      }
      try{
        await axios.post(urls.clients, {
          server: values.server,
          protocol: GRPC,
          options: {
            protos: values.protos.map(file => file.response.filepath),
          },
        })
      } catch (err) {
        return message.error(err.response.data)
      }
      this.props.onSubmit()
    })
  }

  removeFile = async (file) => {
    try {
      await axios.delete(urls.files, {params: {filepath: file.response.filepath}})
    } catch (err) {
      return message.error(err.response.data)
    }
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
