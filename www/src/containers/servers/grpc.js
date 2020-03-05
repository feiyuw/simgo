import React from 'react'
import axios from 'axios'
import { Popconfirm, Divider, Icon, Collapse, Card, List, message, Select, Row, Col, Button, Form, Input } from 'antd'
import './grpc.css'


const {Panel} = Collapse


const MessageItem = ({direction, from, to, ts, body}) => (
  <div className='server-msg-item' direction={direction}>
    <p style={{fontSize: '0.5em'}}>{from} -> {to} ({(new Date(ts)).toLocaleString()})</p>
    <p>{body}</p>
  </div>
)


export default class GrpcServerComponent extends React.Component {
  state = {loadingHandler: true, loadingMessage: true}
  handlers = [{method: 'hellworld.service.SayHello', type: 'raw', content: '{"message": "hello"}'}]
  messages = [
    {direction: 'in', from: '127.0.0.1:12567', to: 'helloworld', ts:1234567890123, body:'{"name": "hello"}'},
    {direction: 'out', from: 'helloworld', to: '127.0.0.1:12567', ts:1234567890125, body:'{"message": "hellohello"}'},
  ]

  async componentDidMount() {
    const {current} = this.props
    if (current !== undefined) {
      await this.loadMessages(current.name)
    }
  }

  loadHandlers = async () => {
    this.setState({loadingHandler: false})
  }

  loadMessages = async (serverName) => {
    let resp
    try {
      resp = await axios.get(`/api/v1/servers/messages?name=${serverName}&skip=0&limit=30`)
    } catch(err) {
      return message.error(err)
    }
    this.messages = resp.data
    this.setState({loadingMessage: false})
  }

  onEditHandler = item => {
    console.log('edit', item)
  }

  onDeleteHandler = item => {
    console.log('delete', item)
  }

  render() {
    const {loadingHandler, loadingMessage} = this.state

    return (
      <div>
        <Collapse onChange={this.loadHandlers}>
          <Panel header='Method Handlers'>
            <List
              size='small'
              dataSource={this.handlers}
              loading={loadingHandler}
              renderItem={item => (
                <List.Item>
                  <List.Item.Meta
                    title={item.method}
                    description={`${item.type}: ${item.content}`}
                  />
                  <Button size='small' type='default' onClick={this.onEditHandler.bind(null, item)}>
                    <Icon type='edit'/>
                  </Button>
                  <Divider type="vertical" />
                  <Popconfirm
                    title='delete this method handler'
                    onConfirm={this.onDeleteHandler.bind(null, item)}
                  >
                    <Button size='small' type='danger'>
                      <Icon type='delete'/>
                    </Button>
                  </Popconfirm>
                </List.Item>
              )}
            />
          </Panel>
        </Collapse>
        <Card title='Messages' style={{marginTop: '15px'}} size='small'>
          <List
            size='small'
            dataSource={this.messages}
            loading={loadingMessage}
            renderItem={msg => (
              <List.Item>
                <MessageItem {...msg} />
              </List.Item>
            )}
          />
        </Card>
      </div>
    )
  }
}
