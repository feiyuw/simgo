import React from 'react'
import axios from 'axios'
import { message, Select, Row, Col, Button, Form, Input } from 'antd'


const {Option} = Select


class GrpcServerComponent extends React.Component {
  state = {loading: true}

  services = []
  methods = []

  response = undefined

  async componentDidMount() {
    await this.fetchServices()
    await this.fetchMethods()
    this.setState({loading: false})
  }

  fetchServices = async () => {
    let resp

    if (this.props.current === undefined) {
      return
    }

    try {
      resp = await axios.get(`/api/v1/grpc/services?clientId=${this.props.current.id}`)
    } catch (err) {
      message.error('fetch grpc services error!')
    }

    this.services = resp.data
    const { setFieldsValue } = this.props.form
    setFieldsValue({service: this.services[0]})
  }

  fetchMethods = async () => {
    let resp

    const { getFieldValue, setFieldsValue } = this.props.form
    const svcName = getFieldValue('service')
    if (this.props.current === undefined || getFieldValue('service') === undefined) {
      return
    }
    try {
      resp = await axios.get(`/api/v1/grpc/methods?clientId=${this.props.current.id}&service=${svcName}`)
    } catch (err) {
      message.error('fetch grpc methods error!')
    }

    this.methods = resp.data
    setFieldsValue({method: this.methods[0]})
  }

  handleSubmit = e => {
    e.preventDefault()
    this.props.form.validateFields(async (err, values) => {
      if (err !== undefined && err !== null) {
        return
      }

      let resp

      try{
        resp = await axios.post('/api/v1/clients/invoke', {
          clientId: this.props.current.id,
          method: values.method,
          data: values.data
        })
      } catch(err) {
        return message.error(err.response.data)
      }

      this.response = resp.data
      this.setState({loading: false})
    })
  }

  handleServiceSwitch = async () => {
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
