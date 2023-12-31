import { useState, useEffect } from 'react'
import { Form, Button, Card, message, Input, Row, Col, Space, List, Tag, Avatar, Popconfirm } from 'antd'
import { EventsOn } from '../wailsjs/runtime'
// @ts-ignore
import { SelectFolder, GetVideoFromPlaylist, StartDownload, StopDownload } from '../wailsjs/go/main/App'
import { IValues, VideoState } from './type'
import { DeleteOutlined } from '@ant-design/icons'
import { LiteralUnion } from 'antd/lib/_util/type'
import { PresetColorType, PresetStatusColorType } from 'antd/es/_util/colors'
import { capitalize } from 'lodash'

const StatusColor: {
  [key in VideoState["Status"]]: LiteralUnion<PresetColorType | PresetStatusColorType>
} = {
  "pending": "default",
  "downloading": "#8950fc",
  "done": "#87d068",
  "error": "#f64e60"
}

const defaultValues: Partial<IValues> = {
  
}

function App(): JSX.Element {
  const [values, setValues] = useState<Partial<IValues>>(defaultValues)
  const [started, setStarted] = useState(false)
  const [loading, setLoading] = useState(false)
  const [loadingPlaylist, setLoadingPlaylist] = useState(false)
  const [videos, setVideos] = useState<VideoState[]>([])
  const onDelete = (videoID: string) => {
    setVideos(videos.filter((video) => video.ID !== videoID))
  }
  const handleSelectFolder = async (name: string) => {
    try {
      const file = await SelectFolder();
      if (file) {
        setValues({
          ...values,
          [name]: file.folder,
        })
      }
    } catch (ex: any) {
      message.error(ex.message);
    }
  }
  const handleGetPlaylist = async () => {
    try {
      if (!values.playlistId) {
        message.error("Vui lòng nhập playlistID!")
        return
      }
      setLoadingPlaylist(true)
      const videos = await GetVideoFromPlaylist(values.playlistId)
      setVideos(videos.map((video) => ({
        ...video,
        Status: "pending"
      })))
    } catch (ex: any) {
      message.error(ex.message);
    } finally {
      setLoadingPlaylist(false)
    }
  }
  const handleStart = async () => {
    try {
      setLoading(true)
      const lists = videos.filter((video) => video.Status !== "done").map((video) => video.ID)
      if (lists.length === 0) {
        message.error("Không có video nào để tải!")
        return
      }
      if (!values.pathFolder) {
        message.error("Vui lòng chọn folder!")
        return
      }
      await StartDownload(lists, values.pathFolder)
    } catch (ex: any) {
      message.error(ex.message);
    } finally {
      setLoading(false)
    }
  }
  const handleStop = async () => {
    try {
      setLoading(true)
      await StopDownload()
    } catch(ex: any) {
      message.error(ex.message);
    }
  }
  useEffect(() => {
    EventsOn("started", () => {
      message.success("Started!")
      setLoading(false)
      setStarted(true)
    })
    EventsOn("stopped", () => {
      message.success("Stopped!")
      setLoading(false)
      setStarted(false)
    })
    EventsOn("status", (msg) => {
      try {
        const json = JSON.parse(msg)
        setVideos((videos) => {
          const index = videos.findIndex((video) => video.ID === json.ID)
          if (index !== -1) {
            videos[index] = {
              ...videos[index],
              Status: json.Status
            }
          }
          return [...videos]
        })
      } catch (ex: any) {
        message.error(ex.message)
      }
    })
    EventsOn("error", (_, data) => {
      message.error(data)
    })
    EventsOn("message", (_, data) => {
      message.info(data)
    })
  }, [])
  return (
    <div>
      <Card style={{
        marginBottom: 16
      }}>
        <Form
          layout="vertical"
          initialValues={defaultValues}
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                label="PlaylistID"
                rules={[{ required: true, message: "Vui lòng chọn nhập playlistID!" }]}
              >
                <Input
                  value={values.playlistId}
                  onChange={(e) => setValues({
                    ...values,
                    playlistId: e.target.value,
                  })}
                  placeholder={"PLKqSEtQQ-HTZk9k6dQfRcspqwxf6QMhBI"}
                />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                label="Path Folder"
              >
                <Space.Compact block>
                  <Input value={values.pathFolder} disabled />
                  <Button type="primary" onClick={() => handleSelectFolder("pathFolder")}>Select Folder</Button>
                </Space.Compact>
              </Form.Item>
            </Col>
          </Row>
          <Space>
            <Button
              type="primary"
              htmlType="submit"
              loading={loadingPlaylist}
              onClick={handleGetPlaylist}
            >
              Get Playlist
            </Button>
            {!started ? (<Button
              type="primary"
              htmlType="submit"
              loading={loading}
              onClick={handleStart}
            >
              Start
            </Button>) : (<Button
              type="primary"
              htmlType="submit"
              loading={loading}
              onClick={handleStop}
            >
              Stop
            </Button>)}
          </Space>
        </Form>
      </Card>
      <Card
        title="Videos"
      >
        <List
          className="demo-loadmore-list"
          loading={loadingPlaylist}
          itemLayout="horizontal"
          dataSource={videos}
          renderItem={(item) => (
            <List.Item
              actions={[<Popconfirm
                title="Chắc chắn xóa?"
                onConfirm={() => onDelete(item.ID)}
                okText="Yes"
                cancelText="No"
              >
                <Button
                  danger
                  icon={<DeleteOutlined />}
                  loading={loading}
                />
              </Popconfirm>]}
            >
              <List.Item.Meta
                avatar={<Avatar src={item.Thumbnails[0]?.URL} />}
                title={<a href={`https://www.youtube.com/watch?v=${item.ID}`} target="_blank">{item.Title}</a>}
                description={item.Author}
              />
              <Tag color={StatusColor[item.Status]}>{capitalize(item.Status)}</Tag>
            </List.Item>
          )}
        />
      </Card>
    </div>
  )
}

export default App