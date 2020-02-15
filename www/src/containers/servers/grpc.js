import React from 'react'
import axios from 'axios'
import { message, Select, Row, Col, Button, Form, Input } from 'antd'


const {Option} = Select


class GrpcServerComponent extends React.Component {
  state = {loading: true}

  response = undefined

  async componentDidMount() {
    await this.fetchServices()
    await this.fetchMethods()
    this.setState({loading: false})
  }

  render() {
    const { getFieldDecorator, getFieldValue } = this.props.form

    return <Form onSubmit={this.handleSubmit}>
          <Row>
            <Col md={12} style={{paddingRight: 20}}>
              <Form.Item>
                {getFieldDecorator('service', {
                  initialValue: undefined,
                  rules: [{ required: true, message: 'Please select service!' }],
                })(
                  <Select placeholder='grpc service' onChange={this.handleServiceSwitch}>
                    {
                      this.services.map((svc, idx) => (
                        <Option key={idx} value={svc}>{svc}</Option>
                      ))
                    }
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
                  <Select placeholder='grpc method'>
                    {
                      this.methods.map((mtd, idx) => (
                        <Option key={idx} value={mtd}>{mtd.substr(getFieldValue('service').length + 1)}</Option>
                      ))
                    }
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
