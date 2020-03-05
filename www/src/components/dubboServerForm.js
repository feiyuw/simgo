import React from 'react'
import { Button, Form, Input } from 'antd'
import {FormItemLayoutWithOutLabel, TwoColumnsFormItemLayout} from '../constants'


class DubboServerForm extends React.Component {
  handleSubmit = e => {
    e.preventDefault()
    this.props.form.validateFields((err, values) => {
      if (!err) {
        console.log('Received values of form: ', values)
      }
    })
  }

  render() {
    const { getFieldDecorator } = this.props.form
    // serializer: hessian2, thrift, jsonp, etc

    return <Form onSubmit={this.handleSubmit}>
          <Form.Item {...TwoColumnsFormItemLayout} label='Server Address'>
            {getFieldDecorator('server', {
              initialValue: '',
              rules: [{ required: true, message: 'Please input server addr!' }],
            })(
              <Input style={{width: 350}} placeholder='type server addr here' />
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


export default Form.create()(DubboServerForm)
