import React from 'react'
import { Select, Row, Col, Button, Form, Input } from 'antd'


const {Option} = Select


class GrpcClientComponent extends React.Component {
  handleSubmit = e => {
    e.preventDefault()
    this.props.form.validateFields((err, values) => {
      if (!err) {
        console.log('Received values of form: ', values)
      }
    })
  }

  handleServiceSwitch = e => {
    console.log(e)
  }

  handleMethodSwitch = e => {
    console.log(e)
  }

  render() {
    const { getFieldDecorator } = this.props.form

    return <Form onSubmit={this.handleSubmit}>
          <Row>
            <Col md={12} style={{paddingRight: 20}}>
              <Form.Item>
                {getFieldDecorator('service', {
                  initialValue: undefined,
                  rules: [{ required: true, message: 'Please select service!' }],
                })(
                  <Select placeholder='grpc service' onChange={this.handleServiceSwitch}>
                    <Option value='helloworld'>helloworld</Option>
                    <Option value='echo'>echo</Option>
                  </Select>
                )}
              </Form.Item>
            </Col>
            <Col md={12}>
              <Form.Item>
                {getFieldDecorator('method', {
                  initialValue: undefined,
                  rules: [{ required: true, message: 'Please select method!' }],
                })(
                  <Select placeholder='grpc method' onChange={this.handleMethodSwitch}>
                    <Option value='UnaryEcho'>UnaryEcho</Option>
                    <Option value='BidiEcho'>BidiEcho</Option>
                  </Select>
                )}
              </Form.Item>
            </Col>
          </Row>
          <Row>
            <Col md={12} style={{paddingRight: 20}}>
              <Form.Item>
                {getFieldDecorator('data', {
                  initialValue: '',
                  rules: [{ required: true, message: 'Request data should be set!' }],
                })(
                  <Input.TextArea rows={20} />
                )}
              </Form.Item>
            </Col>
            <Col md={12}>
              <Input.TextArea rows={20} value='demo response' />
            </Col>
          </Row>
          <Form.Item>
            <Button type="primary" htmlType="submit">
              Connect
            </Button>
          </Form.Item>
        </Form>
  }
}


export default Form.create()(GrpcClientComponent)
