import React from 'react'
import { message, Input, Select, Icon, Button, List, Row, Col } from 'antd'
import axios from 'axios'
import {NewClientDialog} from '../../components/dialog'
import GrpcClientComponent from './grpc'


export default class ClientApp extends React.Component {
  state = {current: undefined, loading: true, showNewDialog: false}
  clients = []
  protocol = 'grpc' // default new client protocol

  async componentDidMount() {
    await this.fetchClients()
    this.setState({current: this.clients[this.clients.length - 1], loading: false})
  }

  fetchClients = async () => {
    let resp
    try {
      resp = await axios.get('/api/v1/clients')
    } catch(err) {
      return message.error(err)
    }

    this.clients = resp.data
  }

  switchProtocol = (protocol) => {
    this.protocol = protocol
  }

  getClientComponent = (protocol) => {
    switch(protocol) {
      case 'grpc':
        return <GrpcClientComponent />
      case 'http':
        return <div>http client</div>
      case 'dubbo':
        return <div>dubbo client</div>
      default:
        return <Icon type='loading' />
    }
  }

  addClient = () => {
    this.setState({showNewDialog: true})
  }

  render() {
    const {current, loading, showNewDialog} = this.state

    return (
      <div>
        <Row>
          <Col md={4}>
            <Input.Group compact>
              <Select defaultValue={this.protocol} style={{width: '60%'}} onChange={this.switchProtocol}>
                <Select.Option value='grpc'>gRPC</Select.Option>
                <Select.Option value='http'>HTTP</Select.Option>
                <Select.Option value='dubbo'>Dubbo</Select.Option>
              </Select>
              <Button style={{width: '40%'}} onClick={this.addClient}>
                <Icon type='plus'/> New
              </Button>
            </Input.Group>
            <List
              size='small'
              dataSource={this.clients}
              bordered={false}
              loading={loading}
              renderItem={item => (
                <List.Item>
                  <Button
                    type='link'
                    style={(current!==undefined && current.id===item.id) ? {backgroundColor: '#1890FF', color: 'white'}: {}}
                    onClick={() => this.setState({current: item})}
                  >
                      # {item.id} {item.protocol} {item.server}
                  </Button>
                </List.Item>
              )}
            />
          </Col>
          <Col md={20} style={{paddingRight: 20, paddingLeft: 15}}>
            {
              current===undefined ?
                <Button style={{width: '100%'}} type='primary' onClick={this.addClient}>
                  <Icon type='plus'/> New Client
                </Button> :
                this.getClientComponent(current.protocol)
            }
          </Col>
        </Row>
        <NewClientDialog
          visible={showNewDialog}
          protocol={this.protocol}
          onClose={() => this.setState({showNewDialog: false})}
          onSubmit={async () => {
            await this.fetchClients()
            this.setState({
              current: this.clients[this.clients.length - 1],
              loading: false,
              showNewDialog: false})
          }}
        />
      </div>
    )
  }
}
