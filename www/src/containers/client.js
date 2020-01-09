import React from 'react'
import { Row, Col, Button, Form, Select, Input } from 'antd'

const { Option } = Select


class clientForm extends React.Component {
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

    return <Form layout="inline" onSubmit={this.handleSubmit}>
          <Form.Item>
            {getFieldDecorator('protocol', {
              initialValue: 'grpc',
              rules: [{ required: true, message: 'Please select one protocol!' }],
            })(
              <Select style={{width: 120}}>
                <Option value='grpc'>gRPC</Option>
                <Option value='http'>HTTP</Option>
                <Option value='dubbo'>Dubbo</Option>
              </Select>
            )}
          </Form.Item>
          <Form.Item>
            {getFieldDecorator('server_addr', {
              initialValue: '',
              rules: [{ required: true, message: 'Please input server addr!' }],
            })(
              <Input style={{width: 350}} placeholder='type server addr here' />
            )}
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit">
              Save
            </Button>
          </Form.Item>
        </Form>
  }
}


const WrapClientForm = Form.create({name: 'client form'})(clientForm)


export default class ClientApp extends React.Component {
  render() {
    return (
      <div>
        <WrapClientForm />
        <Row>
          <Col md={12} style={{paddingRight: 20}}>
            <Input.TextArea rows={10} />
          </Col>
          <Col md={12}>
            <Input.TextArea rows={10} />
          </Col>
        </Row>
        <Row>
          <Col md={24} style={{marginTop: 15}}>
            <Button type="primary" htmlType="submit">
              Send
            </Button>
          </Col>
        </Row>
      </div>
    )
  }
}
