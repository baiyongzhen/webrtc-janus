vagrant@earth:/vagrant/projects/webrtc-janus/cmd/videoroom/rtp-forwarding$ go run main.go -container-path=/vagrant/projects/webrtc-janus/assets/01.mp4
ws://localhost:8188/janus
{
>   "janus": "create",
>   "transaction": "ce4uve8ffnff8i2jum40"
>}
{
<   "janus": "success",
<   "transaction": "ce4uve8ffnff8i2jum40",
<   "data": {
<      "id": 3270870797002896
<   }
<}
{
>   "janus": "attach",
>   "plugin": "janus.plugin.videoroom",
>   "session_id": 3270870797002896,
>   "transaction": "ce4uve8ffnff8i2jum4g"
>}
{
<   "janus": "success",
<   "session_id": 3270870797002896,
<   "transaction": "ce4uve8ffnff8i2jum4g",
<   "data": {
<      "id": 4153612663869414
<   }
<}
{
>   "body": {
>      "display": "test",
>      "ptype": "publisher",
>      "request": "join",
>      "room": 1234
>   },
>   "handle_id": 4153612663869414,
>   "janus": "message",
>   "session_id": 3270870797002896,
>   "transaction": "ce4uve8ffnff8i2jum50"
>}
{
>   "janus": "keepalive",
>   "session_id": 3270870797002896,
>   "transaction": "ce4uve8ffnff8i2jum5g"
>}
{
<   "janus": "ack",
<   "session_id": 3270870797002896,
<   "transaction": "ce4uve8ffnff8i2jum50"
<}
{
<   "janus": "event",
<   "session_id": 3270870797002896,
<   "transaction": "ce4uve8ffnff8i2jum50",
<   "sender": 4153612663869414,
<   "plugindata": {
<      "plugin": "janus.plugin.videoroom",
<      "data": {
<         "videoroom": "joined",
<         "room": 1234,
<         "description": "Demo Room",
<         "id": 8785030944646752,
<         "private_id": 9802063,
<         "publishers": []
<      }
<   }
<}
{
<   "janus": "ack",
<   "session_id": 3270870797002896,
<   "transaction": "ce4uve8ffnff8i2jum5g"
<}
FeedID:= 8785030944646752
{
>   "body": {
>      "audio": true,
>      "bitrate": 128000,
>      "bitrate_cap": true,
>      "data": false,
>      "request": "publish",
>      "video": true
>   },
>   "handle_id": 4153612663869414,
>   "janus": "message",
>   "jsep": {
>      "sdp": "v=0\r\no=- 2773062412132628187 1669984185 IN IP4 0.0.0.0\r\ns=-\r\nt=0 0\r\na=fingerprint:sha-256 E5:AC:B8:07:85:05:AE:01:90:B2:D5:32:9E:58:F6:F5:29:26:55:DB:7E:F7:CC:9D:22:A0:6F:42:FF:28:37:C7\r\na=group:BUNDLE 0 1\r\nm=audio 9 UDP/TLS/RTP/SAVPF 111 9 0 8\r\nc=IN IP4 0.0.0.0\r\na=setup:actpass\r\na=mid:0\r\na=ice-ufrag:eROfthpCTUPYAdKI\r\na=ice-pwd:IwusUVzfdIiVKurLnMmYucqkXsCGIhMv\r\na=rtcp-mux\r\na=rtcp-rsize\r\na=rtpmap:111 opus/48000/2\r\na=fmtp:111 minptime=10;useinbandfec=1\r\na=rtcp-fb:111 transport-cc \r\na=rtpmap:9 G722/8000\r\na=rtcp-fb:9 transport-cc \r\na=rtpmap:0 PCMU/8000\r\na=rtcp-fb:0 transport-cc \r\na=rtpmap:8 PCMA/8000\r\na=rtcp-fb:8 transport-cc \r\na=extmap:1 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01\r\na=ssrc:4115780217 cname:synced-video\r\na=ssrc:4115780217 msid:synced-video synced-video\r\na=ssrc:4115780217 mslabel:synced-video\r\na=ssrc:4115780217 label:synced-video\r\na=msid:synced-video synced-video\r\na=sendrecv\r\na=candidate:3565599043 1 udp 2130706431 10.0.2.15 50287 typ host\r\na=candidate:3565599043 2 udp 2130706431 10.0.2.15 50287 typ host\r\na=candidate:3769770935 1 udp 2130706431 192.168.56.167 40364 typ host\r\na=candidate:3769770935 2 udp 2130706431 192.168.56.167 40364 typ host\r\na=candidate:3528925834 1 udp 2130706431 172.18.0.1 42869 typ host\r\na=candidate:3528925834 2 udp 2130706431 172.18.0.1 42869 typ host\r\na=candidate:233762139 1 udp 2130706431 172.17.0.1 32979 typ host\r\na=candidate:233762139 2 udp 2130706431 172.17.0.1 32979 typ host\r\na=candidate:2724024701 1 udp 1694498815 125.177.129.71 54019 typ srflx raddr 0.0.0.0 rport 58952\r\na=candidate:2724024701 2 udp 1694498815 125.177.129.71 54019 typ srflx raddr 0.0.0.0 rport 58952\r\na=end-of-candidates\r\nm=video 9 UDP/TLS/RTP/SAVPF 96 97 98 99 100 101 102 121 127 120 125 107 108 109 123 118 116\r\nc=IN IP4 0.0.0.0\r\na=setup:actpass\r\na=mid:1\r\na=ice-ufrag:eROfthpCTUPYAdKI\r\na=ice-pwd:IwusUVzfdIiVKurLnMmYucqkXsCGIhMv\r\na=rtcp-mux\r\na=rtcp-rsize\r\na=rtpmap:96 VP8/90000\r\na=rtcp-fb:96 goog-remb \r\na=rtcp-fb:96 ccm fir\r\na=rtcp-fb:96 nack \r\na=rtcp-fb:96 nack pli\r\na=rtcp-fb:96 nack \r\na=rtcp-fb:96 nack pli\r\na=rtcp-fb:96 transport-cc \r\na=rtpmap:97 rtx/90000\r\na=fmtp:97 apt=96\r\na=rtcp-fb:97 nack \r\na=rtcp-fb:97 nack pli\r\na=rtcp-fb:97 transport-cc \r\na=rtpmap:98 VP9/90000\r\na=fmtp:98 profile-id=0\r\na=rtcp-fb:98 goog-remb \r\na=rtcp-fb:98 ccm fir\r\na=rtcp-fb:98 nack \r\na=rtcp-fb:98 nack pli\r\na=rtcp-fb:98 nack \r\na=rtcp-fb:98 nack pli\r\na=rtcp-fb:98 transport-cc \r\na=rtpmap:99 rtx/90000\r\na=fmtp:99 apt=98\r\na=rtcp-fb:99 nack \r\na=rtcp-fb:99 nack pli\r\na=rtcp-fb:99 transport-cc \r\na=rtpmap:100 VP9/90000\r\na=fmtp:100 profile-id=1\r\na=rtcp-fb:100 goog-remb \r\na=rtcp-fb:100 ccm fir\r\na=rtcp-fb:100 nack \r\na=rtcp-fb:100 nack pli\r\na=rtcp-fb:100 nack \r\na=rtcp-fb:100 nack pli\r\na=rtcp-fb:100 transport-cc \r\na=rtpmap:101 rtx/90000\r\na=fmtp:101 apt=100\r\na=rtcp-fb:101 nack \r\na=rtcp-fb:101 nack pli\r\na=rtcp-fb:101 transport-cc \r\na=rtpmap:102 H264/90000\r\na=fmtp:102 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42001f\r\na=rtcp-fb:102 goog-remb \r\na=rtcp-fb:102 ccm fir\r\na=rtcp-fb:102 nack \r\na=rtcp-fb:102 nack pli\r\na=rtcp-fb:102 nack \r\na=rtcp-fb:102 nack pli\r\na=rtcp-fb:102 transport-cc \r\na=rtpmap:121 rtx/90000\r\na=fmtp:121 apt=102\r\na=rtcp-fb:121 nack \r\na=rtcp-fb:121 nack pli\r\na=rtcp-fb:121 transport-cc \r\na=rtpmap:127 H264/90000\r\na=fmtp:127 level-asymmetry-allowed=1;packetization-mode=0;profile-level-id=42001f\r\na=rtcp-fb:127 goog-remb \r\na=rtcp-fb:127 ccm fir\r\na=rtcp-fb:127 nack \r\na=rtcp-fb:127 nack pli\r\na=rtcp-fb:127 nack \r\na=rtcp-fb:127 nack pli\r\na=rtcp-fb:127 transport-cc \r\na=rtpmap:120 rtx/90000\r\na=fmtp:120 apt=127\r\na=rtcp-fb:120 nack \r\na=rtcp-fb:120 nack pli\r\na=rtcp-fb:120 transport-cc \r\na=rtpmap:125 H264/90000\r\na=fmtp:125 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42e01f\r\na=rtcp-fb:125 goog-remb \r\na=rtcp-fb:125 ccm fir\r\na=rtcp-fb:125 nack \r\na=rtcp-fb:125 nack pli\r\na=rtcp-fb:125 nack \r\na=rtcp-fb:125 nack pli\r\na=rtcp-fb:125 transport-cc \r\na=rtpmap:107 rtx/90000\r\na=fmtp:107 apt=125\r\na=rtcp-fb:107 nack \r\na=rtcp-fb:107 nack pli\r\na=rtcp-fb:107 transport-cc \r\na=rtpmap:108 H264/90000\r\na=fmtp:108 level-asymmetry-allowed=1;packetization-mode=0;profile-level-id=42e01f\r\na=rtcp-fb:108 goog-remb \r\na=rtcp-fb:108 ccm fir\r\na=rtcp-fb:108 nack \r\na=rtcp-fb:108 nack pli\r\na=rtcp-fb:108 nack \r\na=rtcp-fb:108 nack pli\r\na=rtcp-fb:108 transport-cc \r\na=rtpmap:109 rtx/90000\r\na=fmtp:109 apt=108\r\na=rtcp-fb:109 nack \r\na=rtcp-fb:109 nack pli\r\na=rtcp-fb:109 transport-cc \r\na=rtpmap:123 H264/90000\r\na=fmtp:123 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=640032\r\na=rtcp-fb:123 goog-remb \r\na=rtcp-fb:123 ccm fir\r\na=rtcp-fb:123 nack \r\na=rtcp-fb:123 nack pli\r\na=rtcp-fb:123 nack \r\na=rtcp-fb:123 nack pli\r\na=rtcp-fb:123 transport-cc \r\na=rtpmap:118 rtx/90000\r\na=fmtp:118 apt=123\r\na=rtcp-fb:118 nack \r\na=rtcp-fb:118 nack pli\r\na=rtcp-fb:118 transport-cc \r\na=rtpmap:116 ulpfec/90000\r\na=rtcp-fb:116 nack \r\na=rtcp-fb:116 nack pli\r\na=rtcp-fb:116 transport-cc \r\na=extmap:1 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01\r\na=ssrc:3288866997 cname:synced-video\r\na=ssrc:3288866997 msid:synced-video synced-video\r\na=ssrc:3288866997 mslabel:synced-video\r\na=ssrc:3288866997 label:synced-video\r\na=msid:synced-video synced-video\r\na=sendrecv\r\n",
>      "trickle": false,
>      "type": "offer"
>   },
>   "session_id": 3270870797002896,
>   "transaction": "ce4uve8ffnff8i2jum60"
>}
{
<   "janus": "ack",
<   "session_id": 3270870797002896,
<   "transaction": "ce4uve8ffnff8i2jum60"
<}
{
<   "janus": "event",
<   "session_id": 3270870797002896,
<   "transaction": "ce4uve8ffnff8i2jum60",
<   "sender": 4153612663869414,
<   "plugindata": {
<      "plugin": "janus.plugin.videoroom",
<      "data": {
<         "videoroom": "event",
<         "room": 1234,
<         "configured": "ok",
<         "audio_codec": "opus",
<         "video_codec": "h264",
<         "streams": [
<            {
<               "type": "audio",
<               "mindex": 0,
<               "mid": "0",
<               "codec": "opus",
<               "fec": true
<            },
<            {
<               "type": "video",
<               "mindex": 1,
<               "mid": "1",
<               "codec": "h264",
<               "h264_profile": "42001f"
<            }
<         ]
<      }
<   },
<   "jsep": {
<      "type": "answer",
<      "sdp": "v=0\r\no=- 2773062412132628187 1669984185 IN IP4 172.18.0.2\r\ns=VideoRoom 1234\r\nt=0 0\r\na=group:BUNDLE 0 1\r\na=ice-options:trickle\r\na=fingerprint:sha-256 06:2E:2E:0A:FC:36:06:F6:29:C6:FB:B2:94:03:1C:14:41:3E:74:EE:35:08:2F:74:A5:81:4C:9D:F2:37:25:1B\r\na=extmap-allow-mixed\r\na=msid-semantic: WMS janus\r\nm=audio 9 UDP/TLS/RTP/SAVPF 111\r\nc=IN IP4 172.18.0.2\r\na=recvonly\r\na=mid:0\r\na=rtcp-mux\r\na=ice-ufrag:Q4b9\r\na=ice-pwd:flEvuSEnQEzVYrK1/nOMP0\r\na=ice-options:trickle\r\na=setup:active\r\na=rtpmap:111 opus/48000/2\r\na=fmtp:111 useinbandfec=1\r\na=msid:janus janus0\r\na=ssrc:3093773747 cname:janus\r\na=candidate:1 1 udp 2015363327 172.18.0.2 44553 typ host\r\na=end-of-candidates\r\nm=video 9 UDP/TLS/RTP/SAVPF 102 121\r\nc=IN IP4 172.18.0.2\r\na=recvonly\r\na=mid:1\r\na=rtcp-mux\r\na=ice-ufrag:Q4b9\r\na=ice-pwd:flEvuSEnQEzVYrK1/nOMP0\r\na=ice-options:trickle\r\na=setup:active\r\na=rtpmap:102 H264/90000\r\na=rtcp-fb:102 ccm fir\r\na=rtcp-fb:102 nack\r\na=rtcp-fb:102 nack pli\r\na=rtcp-fb:102 goog-remb\r\na=rtcp-fb:102 transport-cc\r\na=fmtp:102 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42001f\r\na=extmap:1 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01\r\na=rtpmap:121 rtx/90000\r\na=fmtp:121 apt=102\r\na=msid:janus janus1\r\na=ssrc:1473174889 cname:janus\r\na=candidate:1 1 udp 2015363327 172.18.0.2 44553 typ host\r\na=end-of-candidates\r\n"
<   }
<}
{
>   "body": {
>      "host": "192.168.56.168",
>      "publisher_id": 8785030944646752,
>      "request": "rtp_forward",
>      "room": 1234,
>      "secret": "adminpwd",
>      "streams": [
>         {
>            "mid": "0",
>            "port": 5006
>         },
>         {
>            "mid": "1",
>            "port": 5011
>         }
>      ]
>   },
>   "handle_id": 4153612663869414,
>   "janus": "message",
>   "session_id": 3270870797002896,
>   "transaction": "ce4uve8ffnff8i2jum6g"
>}
Connection State has changed checking
{
<   "janus": "success",
<   "session_id": 3270870797002896,
<   "transaction": "ce4uve8ffnff8i2jum6g",
<   "sender": 4153612663869414,
<   "plugindata": {
<      "plugin": "janus.plugin.videoroom",
<      "data": {
<         "publisher_id": 8785030944646752,
<         "forwarders": [
<            {
<               "stream_id": 4137506855,
<               "host": "192.168.56.168",
<               "port": 5006,
<               "type": "audio"
<            },
<            {
<               "stream_id": 3066185642,
<               "host": "192.168.56.168",
<               "port": 5011,
<               "type": "video"
<            }
<         ],
<         "room": 1234,
<         "videoroom": "rtp_forward"
<      }
<   }
<}
Connection State has changed connected
{
<   "janus": "webrtcup",
<   "session_id": 3270870797002896,
<   "sender": 4153612663869414
<}
2022/12/02 12:29:45 WebRTCUp type  4153612663869414
{
<   "janus": "media",
<   "session_id": 3270870797002896,
<   "sender": 4153612663869414,
<   "mid": "0",
<   "type": "audio",
<   "receiving": true
<}
{
<   "janus": "media",
<   "session_id": 3270870797002896,
<   "sender": 4153612663869414,
<   "mid": "1",
<   "type": "video",
<   "receiving": true
<}
2022/12/02 12:29:45 MediaEvent type video  receiving  true
2022/12/02 12:29:45 MediaEvent type audio  receiving  true
{
>   "janus": "keepalive",
>   "session_id": 3270870797002896,
>   "transaction": "ce4uvfgffnff8i2jum70"
>}
{
<   "janus": "ack",
<   "session_id": 3270870797002896,
<   "transaction": "ce4uvfgffnff8i2jum70"
<}
{
>   "janus": "keepalive",
>   "session_id": 3270870797002896,
>   "transaction": "ce4uvgoffnff8i2jum7g"
>}
{
<   "janus": "ack",
<   "session_id": 3270870797002896,
<   "transaction": "ce4uvgoffnff8i2jum7g"
<}





```base
gst-launch-1.0 \
  audiotestsrc ! \
    audioresample ! audio/x-raw,channels=1,rate=16000 ! \
    opusenc bitrate=20000 ! \
      rtpopuspay ! udpsink host=192.168.56.168 port=5006
```