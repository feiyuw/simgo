import React from 'react'
import axios from 'axios'
import { message, Select, Row, Col, Button, Form, Input } from 'antd'
import urls from '../../urls'


const {Option} = Select


class GrpcClientComponent extends React.Component {
  state = {loading: true}
  clientId = this.props.current && this.props.current.id

  services = []
  methods = []

  response = undefined

  async componentDidMount() {
    await this.fetchServices()
    await this.fetchMethods()
    this.setState({loading: false})
  }

  async componentDidUpdate(prevProps) {
    if (this.props.current !== prevProps.current) {
      this.clientId = this.props.current && this.props.current.id
      await this.fetchServices()
      await this.fetchMethods()
      this.setState({loading: false})
    }
  }

  fetchServices = async () => {
    let resp

    if (this.clientId === undefined) {
      return
    }

    try {
      resp = await axios.get(urls.grpcClientsServices, {params: {clientId: this.clientId}})
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
    if (this.clientId === undefined || getFieldValue('service') === undefined) {
      return
    }
    try {
      resp = await axios.get(urls.grpcClientsMethods, {params: {clientId: this.clientId, service: svcName}})
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
        resp = await axios.post(urls.grpcClientsInvoke, {
          clientId: this.clientId,
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


export default Form.create()(GrpcClientComponent)
