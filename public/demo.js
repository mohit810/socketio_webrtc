let divSelectRoom = document.getElementById("selectRoom")
let inputRoomNumber = document.getElementById("roomNumber")
let signalingContainer = document.getElementById('signalingContainer')
let createSessionButton = document.getElementsByClassName('createSessionButton')
let remoteSessionDescription = document.getElementById('remoteSessionDescription')
let localSessionDescription = document.getElementById('localSessionDescription')
let video1 = document.getElementById('video1')

let roomNumber, encryptedSdp

/* eslint-env browser */
var log = msg => {
  document.getElementById('logs').innerHTML += msg + '<br>'
}

const socket = io()

window.createSession = isPublisher => {
  if (inputRoomNumber.value === '') {
    alert("please enter a room name.")
  } else{
    roomNumber = inputRoomNumber.value
    let pc = new RTCPeerConnection({
      iceServers: [
        {'urls': 'stun:stun.services.mozilla.com'},
        {'urls': 'stun:stun.l.google.com.19302'}
      ]
    })
  pc.oniceconnectionstatechange = e => log(pc.iceConnectionState)
  pc.onicecandidate = event => {
    if (event.candidate === null) {
      encryptedSdp = btoa(JSON.stringify(pc.localDescription))
      localSessionDescription.value = encryptedSdp
      socket.emit("ready", roomNumber)
    }
  }

    socket.emit("createConnection", roomNumber)
    socket.on('created', room => {
      console.log("console log from created socket:",room)
      navigator.mediaDevices.getUserMedia({video: true, audio: false})
          .then(stream => {
            pc.addStream(video1.srcObject = stream)
            pc.createOffer()
                .then(d => {
                  pc.setLocalDescription(d)
                }).catch(log)
          }).catch(log)
    })

    socket.on('joined', room => {
      console.log("console log from joined socket:",room)
      pc.addTransceiver('video')
      pc.createOffer()
          .then(d => pc.setLocalDescription(d))
          .catch(log)

      pc.ontrack = function (event) {
        var el = video1
        el.srcObject = event.streams[0]
        el.autoplay = true
        el.controls = true
      }
    })

    socket.on('ready', () =>{
      var obj = JSON.parse(JSON.stringify({
        "sdp": encryptedSdp,
        "roomName": roomNumber
      }))
      socket.emit("offer",obj)
    })

    socket.on('answer', (event) =>{
      let sd = event
      remoteSessionDescription.value = sd
      if (sd === '') {
        return alert('Session Description must not be empty')
      }

      try {
        pc.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(sd))))
      } catch (e) {
        alert(e)
      }
    })

  let btns = createSessionButton
  for (let i = 0; i < btns.length; i++) {
    btns[i].style = 'display: none'
  }
  divSelectRoom.style = "display: none"
  signalingContainer.style = 'display: block'
}
}
