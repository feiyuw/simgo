import React from 'react'
import axios from 'axios'
import { message, Select, Button, Form, Input } from 'antd'
import {FormItemLayoutWithOutLabel, TwoColumnsFormItemLayout} from '../constants'


const {Option} = Select


class GrpcHandlerForm extends React.Component {
  state = {loading: true}
  methods = []

  async componentDidMount() {
    if (this.props.server === undefined) {
      return
    }

    let resp
    try {
      resp = await axios.get(`/api/v1/servers/grpc/methods?name=${this.props.server}`)
    } catch(err) {
      return message.error(err)
    }
    this.methods = resp.data
    this.setState({loading: false})
  }

  handleSubmit = e => {
    e.preventDefault()
    this.props.form.validateFields(async (err, values) => {
      if (err!==null && err!==undefined) {
        return
      }
      try{
        await axios.post('/api/v1/servers/handlers', {
          name: this.props.server,
          method: values.method,
          type: values.type,
          content: values.content,
        })
      } catch (err) {
        return message.error('failed to add new handler')
      }
      this.props.onSubmit()
    })
  }

  render() {
    const { getFieldDecorator, getFieldValue } = this.props.form

    return <Form onSubmit={this.handleSubmit}>
          <Form.Item {...TwoColumnsFormItemLayout} label='Method'>
            {getFieldDecorator('method', {
              initialValue: undefined,
              rules: [{ required: true, message: 'Please select method!' }],
            })(
              <Select placeholder='grpc method'>
                {
                  this.methods.map((mtd, idx) => (
                    <Option key={idx} value={mtd}>{mtd}</Option>
                  ))
                }
              </Select>
            )}
          </Form.Item>
          <Form.Item {...TwoColumnsFormItemLayout} label='Type'>
            {getFieldDecorator('type', {
              initialValue: 'raw',
            })(
              <Select>
                <Option value='raw'>raw</Option>
                <Option value='javascript'>javascript</Option>
              </Select>
            )}
          </Form.Item>
          <Form.Item {...TwoColumnsFormItemLayout} label='Content'>
            {getFieldDecorator('content', {
              initialValue: '',
              rules: [{ required: true, message: 'content should be set!' }],
            })(
              <Input.TextArea rows={10} allowClear placeholder={getFieldValue('type')==='raw' ? 'response data in JSON format' : 'handler code as JavaScript language'}/>
            )}
          </Form.Item>
          <Form.Item {...FormItemLayoutWithOutLabel}>
            <Button type="primary" htmlType="submit">
              Add
            </Button>
          </Form.Item>
        </Form>
  }
}


export default Form.create()(GrpcHandlerForm)
