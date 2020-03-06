import React from 'react'
import axios from 'axios'
import { Popconfirm, Icon, Collapse, Card, List, message, Button } from 'antd'
import {NewGrpcMethodHandlerDialog} from '../../components/dialog'
import './grpc.css'


const {Panel} = Collapse


const MessageItem = ({method, direction, from, to, ts, body}) => (
  <div className='server-msg-item' direction={direction}>
    <p style={{fontSize: '0.5em'}}>{method}: {from} -> {to} ({(new Date(ts)).toLocaleString()})</p>
    <p>{body}</p>
  </div>
)


export default class GrpcServerComponent extends React.Component {
  state = {loadingHandler: true, loadingMessage: true, showHandlerDialog: false}
  serverName = this.props.current && this.props.current.name
  handlers = []
  messages = []
  skip = 0
  limit = 20

  async componentDidMount() {
    await this.loadMessages()
    await this.loadHandlers()
  }

  async componentWillReceiveProps(nextProps) {
    if (this.props.current !== nextProps.current) {
      this.serverName = nextProps.current.name
      await this.loadMessages()
      await this.loadHandlers()
    }
  }

  loadHandlers = async () => {
    if (this.serverName === undefined) {
      return
    }

    let resp
    try {
      resp = await axios.get(`/api/v1/servers/handlers?name=${this.serverName}`)
    } catch(err) {
      return message.error(err)
    }
    this.handlers = Object.keys(resp.data).map(mtd => (
      {method: mtd, ...resp.data[mtd]}
    ))
    this.setState({loadingHandler: false})
  }

  loadMessages = async () => {
    if (this.serverName === undefined) {
      return
    }

    let resp
    try {
      resp = await axios.get(`/api/v1/servers/messages?name=${this.serverName}&skip=${this.skip}&limit=${this.limit}`)
    } catch(err) {
      return message.error(err)
    }
    this.messages = resp.data
    this.setState({loadingMessage: false})
  }

  onDeleteHandler = async item => {
    try {
      await axios.delete(`/api/v1/servers/handlers?name=${this.props.current.name}&method=${item.method}`)
    } catch(err) {
      return message.error(err)
    }
    await this.loadHandlers()
  }

  onAddHandler = () => {
    this.setState({showHandlerDialog: true})
  }

  render() {
    const {loadingHandler, loadingMessage, showHandlerDialog} = this.state

    return (
      <div>
        <Collapse>
          <Panel header='Method Handlers'>
            <Button size='small' type='primary' onClick={this.onAddHandler}>
              <Icon type='plus'/> New Method Handler
            </Button>
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
            pagination={{
              onChange: async page => {
                this.skip = (page - 1) * this.limit
                await this.loadMessages()
              },
              pageSize: this.limit,
            }}
            renderItem={msg => (
              <List.Item>
                <MessageItem {...msg} />
              </List.Item>
            )}
          />
        </Card>
        <NewGrpcMethodHandlerDialog
          server={this.props.current && this.props.current.name}
          visible={showHandlerDialog}
          onClose={() => this.setState({showHandlerDialog: false})}
          onSubmit={async () => {
            await this.loadHandlers()
            this.setState({showHandlerDialog: false})
          }}
        />
      </div>
    )
  }
}
