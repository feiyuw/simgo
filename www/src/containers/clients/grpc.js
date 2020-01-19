import React from 'react'
import axios from 'axios'
import { message, Select, Row, Col, Button, Form, Input } from 'antd'


const {Option} = Select


class GrpcClientComponent extends React.Component {
  state = {loading: true}

  services = []
  methods = []
  currentService = undefined
  currentMethod = undefined

  componentDidMount() {
    if (this.props.current !== undefined) {
      axios.get(`/api/v1/grpc/services?clientId=${this.props.current.id}`)
        .then(resp => {
          this.services = resp.data
          this.setState({loading: false})
        })
        .catch(err => {
          message.error('fetch services error!')
        })
    }
  }

  handleSubmit = e => {
    e.preventDefault()
    this.props.form.validateFields((err, values) => {
      if (!err) {
        console.log('Received values of form: ', values)
      }
    })
  }

  handleServiceSwitch = async e => {
    this.currentService = e
    if (this.props.current !== undefined && this.currentService !== undefined) {
      this.methods = await axios.get(`/api/v1/grpc/methods?clientId=${this.props.current.id}&service=${this.currentService}`).data
    }
    console.log(this.methods)
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
                  initialValue: this.services[0],
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
