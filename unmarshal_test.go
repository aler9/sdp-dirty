package sdp

import (
	"strconv"
	"testing"
)

const (
	BaseSDP = "v=0\r\n" +
		"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\r\n" +
		"s=SDP Seminar\r\n"

	SessionInformationSDP = BaseSDP +
		"i=A Seminar on the session description protocol\r\n" +
		"t=3034423619 3042462419\r\n"

	// https://tools.ietf.org/html/rfc4566#section-5
	// Parsers SHOULD be tolerant and also accept records terminated
	// with a single newline character.
	SessionInformationSDPLFOnly = "v=0\n" +
		"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\n" +
		"s=SDP Seminar\n" +
		"i=A Seminar on the session description protocol\n" +
		"t=3034423619 3042462419\n"

	// SessionInformationSDPCROnly = "v=0\r" +
	// 	"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\r" +
	// 	"s=SDP Seminar\r"
	// 	"i=A Seminar on the session description protocol\r" +
	// 	"t=3034423619 3042462419\r"

	// Other SDP parsers (e.g. one in VLC media player) allow
	// empty lines.
	SessionInformationSDPExtraCRLF = "v=0\r\n" +
		"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\r\n" +
		"\r\n" +
		"s=SDP Seminar\r\n" +
		"\r\n" +
		"i=A Seminar on the session description protocol\r\n" +
		"\r\n" +
		"t=3034423619 3042462419\r\n" +
		"\r\n"

	URISDP = BaseSDP +
		"u=http://www.example.com/seminars/sdp.pdf\r\n" +
		"t=3034423619 3042462419\r\n"

	EmailAddressSDP = BaseSDP +
		"e=j.doe@example.com (Jane Doe)\r\n" +
		"t=3034423619 3042462419\r\n"

	PhoneNumberSDP = BaseSDP +
		"p=+1 617 555-6011\r\n" +
		"t=3034423619 3042462419\r\n"

	SessionConnectionInformationSDP = BaseSDP +
		"c=IN IP4 224.2.17.12/127\r\n" +
		"t=3034423619 3042462419\r\n"

	SessionBandwidthSDP = BaseSDP +
		"b=X-YZ:128\r\n" +
		"b=AS:12345\r\n" +
		"t=3034423619 3042462419\r\n"

	TimingSDP = BaseSDP +
		"t=2873397496 2873404696\r\n"

	// Short hand time notation is converted into NTP timestamp format in
	// seconds. Because of that unittest comparisons will fail as the same time
	// will be expressed in different units.
	RepeatTimesSDP = TimingSDP +
		"r=604800 3600 0 90000\r\n" +
		"r=3d 2h 0 21h\r\n"

	RepeatTimesSDPExpected = TimingSDP +
		"r=604800 3600 0 90000\r\n" +
		"r=259200 7200 0 75600\r\n"

	RepeatTimesSDPExtraCRLF = RepeatTimesSDPExpected +
		"\r\n"

	// The expected value looks a bit different for the same reason as mentioned
	// above regarding RepeatTimes.
	TimeZonesSDP = TimingSDP +
		"r=2882844526 -1h 2898848070 0\r\n"

	TimeZonesSDPExpected = TimingSDP +
		"r=2882844526 -3600 2898848070 0\r\n"

	TimeZonesSDP2 = TimingSDP +
		"z=2882844526 -3600 2898848070 0\r\n"

	TimeZonesSDP2ExtraCRLF = TimeZonesSDP2 +
		"\r\n"

	SessionEncryptionKeySDP = TimingSDP +
		"k=prompt\r\n"

	SessionEncryptionKeySDPExtraCRLF = SessionEncryptionKeySDP +
		"\r\n"

	SessionAttributesSDP = TimingSDP +
		"a=rtpmap:96 opus/48000\r\n"

	MediaNameSDP = TimingSDP +
		"m=video 51372 RTP/AVP 99\r\n" +
		"m=audio 54400 RTP/SAVPF 0 96\r\n"

	MediaNameSDPExtraCRLF = MediaNameSDP +
		"\r\n"

	MediaTitleSDP = MediaNameSDP +
		"i=Vivamus a posuere nisl\r\n"

	MediaConnectionInformationSDP = MediaNameSDP +
		"c=IN IP4 203.0.113.1\r\n"

	MediaConnectionInformationSDPExtraCRLF = MediaConnectionInformationSDP +
		"\r\n"

	MediaDescriptionOutOfOrderSDP = MediaNameSDP +
		"a=rtpmap:99 h263-1998/90000\r\n" +
		"a=candidate:0 1 UDP 2113667327 203.0.113.1 54400 typ host\r\n" +
		"c=IN IP4 203.0.113.1\r\n" +
		"i=Vivamus a posuere nisl\r\n"

	MediaDescriptionOutOfOrderSDPActual = MediaNameSDP +
		"i=Vivamus a posuere nisl\r\n" +
		"c=IN IP4 203.0.113.1\r\n" +
		"a=rtpmap:99 h263-1998/90000\r\n" +
		"a=candidate:0 1 UDP 2113667327 203.0.113.1 54400 typ host\r\n"

	MediaBandwidthSDP = MediaNameSDP +
		"b=X-YZ:128\r\n" +
		"b=AS:12345\r\n"

	MediaEncryptionKeySDP = MediaNameSDP +
		"k=prompt\r\n"

	MediaEncryptionKeySDPExtraCRLF = MediaEncryptionKeySDP +
		"\r\n"

	MediaAttributesSDP = MediaNameSDP +
		"a=rtpmap:99 h263-1998/90000\r\n" +
		"a=candidate:0 1 UDP 2113667327 203.0.113.1 54400 typ host\r\n" +
		"a=rtcp-fb:97 ccm fir\r\n" +
		"a=rtcp-fb:97 nack\r\n" +
		"a=rtcp-fb:97 nack pli\r\n"

	CanonicalUnmarshalSDP = "v=0\r\n" +
		"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\r\n" +
		"s=SDP Seminar\r\n" +
		"i=A Seminar on the session description protocol\r\n" +
		"u=http://www.example.com/seminars/sdp.pdf\r\n" +
		"e=j.doe@example.com (Jane Doe)\r\n" +
		"p=+1 617 555-6011\r\n" +
		"c=IN IP4 224.2.17.12/127\r\n" +
		"b=X-YZ:128\r\n" +
		"b=AS:12345\r\n" +
		"t=2873397496 2873404696\r\n" +
		"t=3034423619 3042462419\r\n" +
		"r=604800 3600 0 90000\r\n" +
		"z=2882844526 -3600 2898848070 0\r\n" +
		"k=prompt\r\n" +
		"a=candidate:0 1 UDP 2113667327 203.0.113.1 54400 typ host\r\n" +
		"a=recvonly\r\n" +
		"m=audio 49170 RTP/AVP 0\r\n" +
		"i=Vivamus a posuere nisl\r\n" +
		"c=IN IP4 203.0.113.1\r\n" +
		"b=X-YZ:128\r\n" +
		"k=prompt\r\n" +
		"a=sendrecv\r\n" +
		"m=video 51372 RTP/AVP 99\r\n" +
		"a=rtpmap:99 h263-1998/90000\r\n"

	SDPWithoutTiming = "v=0\r\n" +
		"m=video 0 RTP/AVP/TCP 96\r\n" +
		"a=rtpmap:96 H265/90000\r\n" +
		"a=fmtp:96 sprop-vps=QAEMAf//AWAAAAMAsAAAAwAAAwB4FwJA; sprop-sps=QgEBAWAAAAMAsAAAAwAAAwB4oAKggC8c1YgXuRZFL/y5/E/qbgQEBAE=; sprop-pps=RAHAcvBTJA==;\r\n" +
		"a=control:streamid=0\r\n" +
		"m=audio 0 RTP/AVP/TCP 97\r\n" +
		"a=rtpmap:97 mpeg4-generic/44100/2\r\n" +
		"a=fmtp:97 profile-level-id=1;mode=AAC-hbr;sizelength=13;indexlength=3;indexdeltalength=3;config=1210\r\n" +
		"a=control:streamid=1\r\n"

	SDPWithoutTiming2 = "v=0\r\n" +
		"s=SDP Seminar\r\n" +
		"m=video 0 RTP/AVP/TCP 96\r\n" +
		"a=rtpmap:96 H265/90000\r\n" +
		"a=fmtp:96 sprop-vps=QAEMAf//AWAAAAMAsAAAAwAAAwB4FwJA; sprop-sps=QgEBAWAAAAMAsAAAAwAAAwB4oAKggC8c1YgXuRZFL/y5/E/qbgQEBAE=; sprop-pps=RAHAcvBTJA==;\r\n" +
		"a=control:streamid=0\r\n" +
		"m=audio 0 RTP/AVP/TCP 97\r\n" +
		"a=rtpmap:97 mpeg4-generic/44100/2\r\n" +
		"a=fmtp:97 profile-level-id=1;mode=AAC-hbr;sizelength=13;indexlength=3;indexdeltalength=3;config=1210\r\n" +
		"a=control:streamid=1\r\n"

	SDPWithoutTiming3 = "v=0\r\n" +
		"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\r\n" +
		"m=video 0 RTP/AVP/TCP 96\r\n" +
		"a=rtpmap:96 H265/90000\r\n" +
		"a=fmtp:96 sprop-vps=QAEMAf//AWAAAAMAsAAAAwAAAwB4FwJA; sprop-sps=QgEBAWAAAAMAsAAAAwAAAwB4oAKggC8c1YgXuRZFL/y5/E/qbgQEBAE=; sprop-pps=RAHAcvBTJA==;\r\n" +
		"a=control:streamid=0\r\n" +
		"m=audio 0 RTP/AVP/TCP 97\r\n" +
		"a=rtpmap:97 mpeg4-generic/44100/2\r\n" +
		"a=fmtp:97 profile-level-id=1;mode=AAC-hbr;sizelength=13;indexlength=3;indexdeltalength=3;config=1210\r\n" +
		"a=control:streamid=1\r\n"

	SDPWithoutTiming4 = "v=0\r\n" +
		"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\r\n" +
		"s=SDP Seminar\r\n" +
		"m=video 0 RTP/AVP/TCP 96\r\n" +
		"a=rtpmap:96 H265/90000\r\n" +
		"a=fmtp:96 sprop-vps=QAEMAf//AWAAAAMAsAAAAwAAAwB4FwJA; sprop-sps=QgEBAWAAAAMAsAAAAwAAAwB4oAKggC8c1YgXuRZFL/y5/E/qbgQEBAE=; sprop-pps=RAHAcvBTJA==;\r\n" +
		"a=control:streamid=0\r\n" +
		"m=audio 0 RTP/AVP/TCP 97\r\n" +
		"a=rtpmap:97 mpeg4-generic/44100/2\r\n" +
		"a=fmtp:97 profile-level-id=1;mode=AAC-hbr;sizelength=13;indexlength=3;indexdeltalength=3;config=1210\r\n" +
		"a=control:streamid=1\r\n"

	SDPLiveReporter = "v=0\r\n" +
		"o=-0 0 IN IP4 127.0.0.1\r\n" +
		"s=No Name\r\n" +
		"c=IN IP4 0.0.0.0\r\n" +
		"t=0 0\r\n" +
		"a=control:*\r\n" +
		"m=video 0 RTP/AVP 96\r\n" +
		"b=AS:253\r\n" +
		"a=rtpmap:96 H264/90000\r\n" +
		"a=fmtp:96 packetization-mode=1; sprop-parameter-sets=J2QAHqxWgKA9pqAgIMBA,KO48sA==; profile-level-id=64001E\r\n" +
		"a=control:streamid=0\r\n" +
		"m=audio 0 RTP/AVP 97\r\n" +
		"b=AS:189\r\n" +
		"a=rtpmap:97 MPEG4-GENERIC/48000/1\r\n" +
		"a=fmtp:97 profile-level-id=1;mode=AAC-hbr;sizelength=13;indexLength=3;indexDeltaLength=3;config=118856E500\r\n" +
		"a=control:streamid=1\r\n"

	SDPSony = "v=0\r\n" +
		"o=- 12345 1 IN IP4 10.21.61.139\r\n" +
		"s=Sony RTSP Server\r\n" +
		"t=0 0\r\n" +
		"a=range:npt=now-\r\n" +
		"c=IN IP4 0.0.0.0\r\n" +
		"m=video 0 RTP/AVP 105\r\n" +
		"a=rtpmap:105 H264/90000\r\n" +
		"a=control:trackID=1\r\n" +
		"a=recvonly\r\n" +
		"a=framerate:25.0\r\n" +
		"a=fmtp:105 packetization-mode=1; profile-level-id=640028; sprop-parameter-sets=Z2QAKKwa0A8ARPy4CIAAAAMAgAAADLWgAtwAHJ173CPFCKg=,KO4ESSJAAAAAAAAAAA==\r\n"

	SDPWithoutTimingButGlobalAttrs = "v=0\r\n" +
		"o=- 1601628642593 1 IN IP4 172.20.1.150\r\n" +
		"s=Session\r\n" +
		"i=Camera arhive stream\r\n" +
		"c=IN IP4 5.19.245.3\r\n" +
		"a=control:rtsp://172.20.1.150:3297/streams/videoa=recvonly\r\n" +
		"a=range:clock=20200908T132728Z-20200909T030844Z;clock=20200909T031019Z-20200909T103133Z;clock=20200909T103300Z-20200909T141237Z\r\n" +
		"a=packetization-supported:DH\r\n" +
		"m=video 0 RTP/AVP 96\r\n" +
		"a=proto:RTP/AVP;unicast\r\n" +
		"a=tool:Archive controller\r\n" +
		"a=rtpmap:96 H264/90000\r\n" +
		"a=fmtp:96 packetization-mode=1;sprop-parameter-sets=Z00AHpY1QWh7TcBAQECA,aO48gA==\r\n"

	SdpWithUnorderedGlobalAttrs = "v=0\r\n" +
		"o=- 1646532490 1646532490 IN IP4 10.175.31.17\r\n" +
		"a=control:*\r\n" +
		"a=source-filter: incl IN IP4 * 10.175.31.17\r\n" +
		"s=RTSP Server\r\n" +
		"a=range:npt=0-\r\n" +
		"t=0 0\r\n" +
		"m=video 0 RTP/AVP 96\r\n" +
		"c=IN IP4 0.0.0.0\r\n" +
		"a=rtpmap:96 H264/90000\r\n" +
		"a=fmtp:96 profile-level-id=4D001E; packetization-mode=1; sprop-parameter-sets=Z00AHpWoKAv+VA==,aO48gA==\r\n" +
		"a=control:?ctype=video\r\n" +
		"a=recvonly\r\n" +
		"m=application 0 RTP/AVP 106\r\n" +
		"a=rtpmap:106 vnd.onvif.metadata/90000\r\n" +
		"a=control:?ctype=app106\r\n" +
		"a=sendonly\r\n"
)

func TestRoundTrip(t *testing.T) {
	for _, test := range []struct {
		Name   string
		SDP    string
		Actual string
	}{
		{
			Name:   "SessionInformationSDPLFOnly",
			SDP:    SessionInformationSDPLFOnly,
			Actual: SessionInformationSDP,
		},
		// {
		// 	Name:   "SessionInformationSDPCROnly",
		// 	SDP:    SessionInformationSDPCROnly,
		// 	Actual: SessionInformationSDPBaseSDP,
		// },
		{
			Name:   "SessionInformationSDPExtraCRLF",
			SDP:    SessionInformationSDPExtraCRLF,
			Actual: SessionInformationSDP,
		},
		{
			Name: "SessionInformation",
			SDP:  SessionInformationSDP,
		},
		{
			Name: "URI",
			SDP:  URISDP,
		},
		{
			Name: "EmailAddress",
			SDP:  EmailAddressSDP,
		},
		{
			Name: "PhoneNumber",
			SDP:  PhoneNumberSDP,
		},
		{
			Name:   "RepeatTimesSDPExtraCRLF",
			SDP:    RepeatTimesSDPExtraCRLF,
			Actual: RepeatTimesSDPExpected,
		},
		{
			Name: "SessionConnectionInformation",
			SDP:  SessionConnectionInformationSDP,
		},
		{
			Name: "SessionBandwidth",
			SDP:  SessionBandwidthSDP,
		},
		{
			Name: "SessionEncryptionKey",
			SDP:  SessionEncryptionKeySDP,
		},
		{
			Name:   "SessionEncryptionKeyExtraCRLF",
			SDP:    SessionEncryptionKeySDPExtraCRLF,
			Actual: SessionEncryptionKeySDP,
		},
		{
			Name: "SessionAttributes",
			SDP:  SessionAttributesSDP,
		},
		{
			Name:   "TimeZonesSDP2ExtraCRLF",
			SDP:    TimeZonesSDP2ExtraCRLF,
			Actual: TimeZonesSDP2,
		},
		{
			Name: "MediaName",
			SDP:  MediaNameSDP,
		},
		{
			Name:   "MediaNameExtraCRLF",
			SDP:    MediaNameSDPExtraCRLF,
			Actual: MediaNameSDP,
		},
		{
			Name: "MediaTitle",
			SDP:  MediaTitleSDP,
		},
		{
			Name: "MediaConnectionInformation",
			SDP:  MediaConnectionInformationSDP,
		},
		{
			Name:   "MediaConnectionInformationExtraCRLF",
			SDP:    MediaConnectionInformationSDPExtraCRLF,
			Actual: MediaConnectionInformationSDP,
		},
		{
			Name:   "MediaDescriptionOutOfOrder",
			SDP:    MediaDescriptionOutOfOrderSDP,
			Actual: MediaDescriptionOutOfOrderSDPActual,
		},
		{
			Name: "MediaBandwidth",
			SDP:  MediaBandwidthSDP,
		},
		{
			Name: "MediaEncryptionKey",
			SDP:  MediaEncryptionKeySDP,
		},
		{
			Name:   "MediaEncryptionKeyExtraCRLF",
			SDP:    MediaEncryptionKeySDPExtraCRLF,
			Actual: MediaEncryptionKeySDP,
		},
		{
			Name: "MediaAttributes",
			SDP:  MediaAttributesSDP,
		},
		{
			Name: "CanonicalUnmarshal",
			SDP:  CanonicalUnmarshalSDP,
		},
		{
			Name: "SDPWithoutTiming",
			SDP:  SDPWithoutTiming,
		},
		{
			Name: "SDPWithoutTiming2",
			SDP:  SDPWithoutTiming2,
		},
		{
			Name: "SDPWithoutTiming3",
			SDP:  SDPWithoutTiming3,
		},
		{
			Name: "SDPWithoutTiming4",
			SDP:  SDPWithoutTiming4,
		},
		{
			Name: "SDPWithoutTimingButGlobalAttrs",
			SDP:  SDPWithoutTimingButGlobalAttrs,
		},
	} {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			sd := &SessionDescription{}

			err := sd.Unmarshal([]byte(test.SDP))
			if got, want := err, error(nil); got != want {
				t.Fatalf("Unmarshal:\nerr=%v\nwant=%v", got, want)
			}

			actual, err := sd.Marshal()
			if got, want := err, error(nil); got != want {
				t.Fatalf("Marshal:\nerr=%v\nwant=%v", got, want)
			}
			want := test.SDP
			if test.Actual != "" {
				want = test.Actual
			}
			if got := string(actual); got != want {
				t.Fatalf("Marshal:\ngot=%s\nwant=%s",
					strconv.Quote(got), strconv.Quote(want),
				)
			}
		})
	}
}

func TestUnmarshalRepeatTimes(t *testing.T) {
	sd := &SessionDescription{}
	if err := sd.Unmarshal([]byte(RepeatTimesSDP)); err != nil {
		t.Errorf("error: %v", err)
	}

	actual, err := sd.Marshal()
	if got, want := err, error(nil); got != want {
		t.Fatalf("Marshal(): err=%v, want %v", got, want)
	}
	if string(actual) != RepeatTimesSDPExpected {
		t.Errorf("error:\n\nEXPECTED:\n%v\nACTUAL:\n%v", RepeatTimesSDPExpected, actual)
	}
}

func TestUnmarshalTimeZones(t *testing.T) {
	sd := &SessionDescription{}
	if err := sd.Unmarshal([]byte(TimeZonesSDP)); err != nil {
		t.Errorf("error: %v", err)
	}

	actual, err := sd.Marshal()
	if got, want := err, error(nil); got != want {
		t.Fatalf("Marshal(): err=%v, want %v", got, want)
	}
	if string(actual) != TimeZonesSDPExpected {
		t.Errorf("error:\n\nEXPECTED:\n%v\nACTUAL:\n%v", TimeZonesSDPExpected, actual)
	}
}

func TestUnmarshalNonNilAddress(t *testing.T) {
	in := "v=0\r\no=0 0 0 IN IP4 0\r\ns=0\r\nc=IN IP4\r\nt=0 0\r\n"
	var sd SessionDescription
	err := sd.Unmarshal([]byte(in))
	if err != nil {
		t.Fatalf("failed to unmarshal %q", in)
	}
	out, err := sd.Marshal()
	if err != nil {
		t.Errorf("failed to marshal unmarshalled %q", in)
	}
	if string(out) != in {
		t.Errorf("round trip = %q want %q", out, in)
	}
}

func TestUnmarshalNonStandard(t *testing.T) {
	for _, test := range []struct {
		Name string
		SDP  string
	}{
		{
			Name: "SDPLiveReporter",
			SDP:  SDPLiveReporter,
		},
		{
			Name: "SDPSony",
			SDP:  SDPSony,
		},
		{
			Name: "SdpWithUnorderedGlobalAttrs",
			SDP:  SdpWithUnorderedGlobalAttrs,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			var sd SessionDescription
			err := sd.Unmarshal([]byte(test.SDP))
			if err != nil {
				t.Fatalf("failed")
			}
		})
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	raw := []byte(CanonicalUnmarshalSDP)
	for i := 0; i < b.N; i++ {
		var sd SessionDescription
		err := sd.Unmarshal(raw)
		if err != nil {
			b.Fatal(err)
		}
	}
}
