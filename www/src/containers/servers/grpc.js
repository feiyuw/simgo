import React from 'react'
import axios from 'axios'
import { message, Select, Row, Col, Button, Form, Input } from 'antd'


const {Option} = Select


class GrpcServerComponent extends React.Component {
  state = {loading: true}

  response = undefined

  render() {
    const { getFieldDecorator, getFieldValue } = this.props.form

    return <Form onSubmit={this.handleSubmit}>
          <Row>
            <Col md={12} style={{paddingRight: 20}}>
              <Form.Item>
                {getFieldDecorator('data', {
                  initialValue: '',
                  rules: [{ required: true, message: 'Request data should be set!' }],
                })(
                  <Input.TextArea rows={20} allowClear placeholder='JSON request data'/>
                )}
              </Form.Item>
            </Col>
            <Col md={12}>
              <Input.TextArea rows={20} value={this.response && JSON.stringify(this.response)} disabled/>
            </Col>
          </Row>
          <Form.Item>
            <Button type="primary" htmlType="submit">
              Send
            </Button>
          </Form.Item>
        </Form>
  }
}


export default Form.create()(GrpcServerComponent)
