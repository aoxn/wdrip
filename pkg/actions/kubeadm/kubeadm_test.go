package kubeadm

import (
	"bytes"
	"fmt"
	"github.com/aoxn/wdrip/pkg/actions"
	v12 "github.com/aoxn/wdrip/pkg/apis/alibabacloud.com/v1"
	"github.com/ghodss/yaml"
	"html/template"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

var bootcfg = `
cloudType: public
clusterid: kubernetes-clusterid-demo
runtime:
  name: runtime
  version: 18.09.2
endpoint:
  intranet: 192.168.0.216
etcd:
  initToken: 9521145d-f16f-4112-add4-aab07c86b7d8
  name: etcd
  peerCA:
    cert: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURHakNDQWdLZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREErTVNjd0ZBWURWUVFLRXcxaGJHbGkKWVdKaElHTnNiM1ZrTUE4R0ExVUVDaE1JYUdGdVozcG9iM1V4RXpBUkJnTlZCQU1UQ210MVltVnlibVYwWlhNdwpIaGNOTVRrd09ERXhNRGt6TmpFNFdoY05Namt3T0RBNE1Ea3pOakU0V2pBK01TY3dGQVlEVlFRS0V3MWhiR2xpCllXSmhJR05zYjNWa01BOEdBMVVFQ2hNSWFHRnVaM3BvYjNVeEV6QVJCZ05WQkFNVENtdDFZbVZ5Ym1WMFpYTXcKZ2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRQy8zampXS0dJSnJLcUxPL1lxakFJSgpTd0pCSVJremwrSE9HL2ZUU2RkaW9obzhGaHgrcDZ0ZzRra0pJU1VsejEyb3Q5Ty9uSk8zcTliZlY4eHVwZkcyCm0xRmNwUDVyeW1sWHBXYnJlSWNoYStwRXpTZXIwSkJWU2tjcnI0bG4yc3lsRjlYOEdBZTVkZS9qY2w4QnhsWXUKNXBHdk93TzdNK2E2QXFPaitWT1lYUXM1d0pUMklFaHhlNHIzTThiYmNUOGtTNHIySlJmTjBWc2RWN2NQV1pubApCRWdmclZ1Y2UxS0M3MitxUUZOS1ZnS3ZnbER0MmpUYzF1UFl0cW8wVWo0MmIzV2lGRTExbFJ2cU5INjN3WjBHCkFiZlBWc0Y2OGZNditSVWdsRUZQbjRDTEIzaHFtUWE3QTI4Q21ocURYcDkwSTZEVEhhSWV1UGJVVlFKMjhMakoKQWdNQkFBR2pJekFoTUE0R0ExVWREd0VCL3dRRUF3SUNwREFQQmdOVkhSTUJBZjhFQlRBREFRSC9NQTBHQ1NxRwpTSWIzRFFFQkN3VUFBNElCQVFBVzEyemJ5T3lUVlFHRkJuTi80SzZRTUNKU25aY2V6cHVoSGFac2F6ZGU2cFF5CnJhNjZVTlBJdVlRNThyK0t4N0tvWGFZRFEwWEFzcDNoT1Zpd3NRWitLSTY5alZoUDhNRmpWQ2N6WHFTM0R5SkIKb2R0OHVON2NUQU42Znk3clEreGtJRFNDaEtza1VzN3dGeE9WeFhVMHdWUE5qREFxTmVJeCtyUHQ2aEdKMDA0dQp2cXhrUnBZa082VTR6R0htdTZFekxEWkJlNDk1Q3pFNEdCcVhVbW9pLzN6MTFmekdYVVY5anZNSTZNaElPT1VlCkV2b3pTNHFyOUZHWVZKMXBPR1VhalovNFVnT0Y1VmF5alNINEo3UzRWOUw0RDM5ZkErVGVkdTRJVVdFcUo1TG0KeWJ1Y3QvRWNpZkx6LzlyN3hmTHRhdDkxYUtHYVdTMzJuYVVxT3BKTwotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    key: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBdjk0NDFpaGlDYXlxaXp2Mktvd0NDVXNDUVNFWk01Zmh6aHYzMDBuWFlxSWFQQlljCmZxZXJZT0pKQ1NFbEpjOWRxTGZUdjV5VHQ2dlczMWZNYnFYeHRwdFJYS1QrYThwcFY2Vm02M2lISVd2cVJNMG4KcTlDUVZVcEhLNitKWjlyTXBSZlYvQmdIdVhYdjQzSmZBY1pXTHVhUnJ6c0R1elBtdWdLam8vbFRtRjBMT2NDVQo5aUJJY1h1Szl6UEcyM0UvSkV1SzlpVVh6ZEZiSFZlM0QxbVo1UVJJSDYxYm5IdFNndTl2cWtCVFNsWUNyNEpRCjdkbzAzTmJqMkxhcU5GSStObTkxb2hSTmRaVWI2alIrdDhHZEJnRzN6MWJCZXZIekwva1ZJSlJCVDUrQWl3ZDQKYXBrR3V3TnZBcG9hZzE2ZmRDT2cweDJpSHJqMjFGVUNkdkM0eVFJREFRQUJBb0lCQVFDK0tnZFZJd05BS1hXQwp1SFJjYVJYZmxndHU5OW9kaTZ5TzlxTmpNKzJZNGFkMDVFbHJzczBtSWtGWEhoWE9hcitlYUV3anZwR2QybUFHClR1UGN5dlpPRVpUTGFQQ05iem1IVi9Vdzd1Mm56MmlLdG5kYVVFV3Rjd2dsSVQ3and6VlBiOWR6bTNHVWZISzkKa0c3ZnVHOVUzc3VIek1yKzhRcitVMzFUR051a3h0K1dIZE5IMGNqam5FOTRTNjB0dVdxbUgxQ3dYZzVGa1hxZQpyeFlKdlZWWVQyaFhpOTJtWVFDTEhTM3RFazJCN1VkWStXYXNHLzBRSVI0Vk1uQ2QxTXFjSWFrZU5ZZGlEQXg1Ci9VeFdqeXB1RWc0d1RZMnYyelBzbXRJY0tFL3pkQUtmeUUzS0xtMzBsUlM4YXdMNUQ4RGJUK2JLMlJTbzNYTEUKUmpyUkdOQWRBb0dCQVA4dW9seTBOTWVBZ2N0elYwQThnQzd5ZnB2YWNxU1JBYWJtUlRrb0w4SlpJcVl3SGRTZwpFemw1NmtyWk5wd1VrVnU5TFFKZXRuN084RXNHeENNaFh2cEpqUlpSNUE1alhwZVRpdzBBeHJZdjgyS1pzdVRsCmVKL1Z0ZEpjMmNwU20yejR3eWYveXl6R2RGMitCTCtoamoxSzFUVWFNY2JraVlFeEloc3lPUnhmQW9HQkFNQjcKcEM3bEJwekJ0VVlIaWpwZWZDNk9LcGplU2xkUjA2WDZtWXF1Vm9EU2FmU3NZZHNqSlYycHAvMWNQV2p6N25LSApIVGlmMlBoQkJuSGpSbFJ6eXFSSnFLZ1RkTVA4U2xUSXNubmlxb3Q1ZUxsb3Iva1V2VVRGWEpTRFI4YVc2R0NoClhLVVFiWVV4TEpsNU5Rbmw4Z2lkQWtEQXBvWmVweDJSbWZnV2tEdlhBb0dCQU9tZmFFWTNOWnJ3cStQMzFRbWEKV2tDaEFnanVGY2RVZW82eWd1MnhQUWhSVXNlVGhid0VVWlZ0YUhMZUtvRDYwNW1KdUl0UzZ3RnRzOUQ4Z05VbwoyQ3VNNnY1a09zWnhjMGlTYTl1YnVsRlIxU0dRVlpmNS9sVlc4ditFd2wvWkFUTGpETWl5QnZFWFA3SVRKWVhNCkFzMWFsWmZvUTJvMEVTK1dMQU42RjQ5QkFvR0FQQWtVYjQvOW5QMEtKanFKMGFUUXhOQ1ExcmRXcHArZURRSDUKeS9pT2dJV2dpTEVQb0lMNHo1cndDNlV1ZmtLL0Y3ZXUxSTkrNUFlY0UvK1lXeFQybW9GaStuRi9GUFhtMVVUMwo3ZTVWMVVUZzg2dUVYNE1wZVg2NVhwVUgyUmdPaUwwcm9VeGJiSlNtM1lPaG1HSEJPUkNIdUZ5ZVZBREh6UUF4Cjd1QlN3bWNDZ1lBcnkvSEhBVnlONThNZ0lmMmw2M2xuNUpRZVVDUWRESWVZTUtqM2o3emoxOU1LcXJwMjg4OGcKYis1TWczeGJVaEk1NGppNENWZ2ZmZERGTXRzbnRlaHp3UUxRRzVkWGZYOE1OK0pYL0w3VnlUaEdIN0JtUEVCVQp0ckNMSXdyc0krUFI0dmRnbWZrNHAyQjZtSVp1aHFwMGdlTlFJUzZvcWdhbUpLNUlVVGJKbFE9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
  serverCA:
    cert: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURHakNDQWdLZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREErTVNjd0ZBWURWUVFLRXcxaGJHbGkKWVdKaElHTnNiM1ZrTUE4R0ExVUVDaE1JYUdGdVozcG9iM1V4RXpBUkJnTlZCQU1UQ210MVltVnlibVYwWlhNdwpIaGNOTVRrd09ERXhNRGt6TmpFNFdoY05Namt3T0RBNE1Ea3pOakU0V2pBK01TY3dGQVlEVlFRS0V3MWhiR2xpCllXSmhJR05zYjNWa01BOEdBMVVFQ2hNSWFHRnVaM3BvYjNVeEV6QVJCZ05WQkFNVENtdDFZbVZ5Ym1WMFpYTXcKZ2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRQ2JhRFB2cXQ2ZGtzQiswU21oUHFrTwpPQnVKeFFNaWRSeEtUWVh2WENtcFNHM0tVc0lNNzc0Y2dpTEx4TWJNbitDNXhkZ3FhVHpDb1pQQjZiSzY2azVJCnd5bm94Q1N2REdGZWxnc0ZPKzdqQU1NR3FhWUVyZFRuL2JDd2M1RTJSNmlaWFRvbjFkSGxsejUxQTkvSlVBY3YKT2FsM1J4L1pZclViaFJaZDg0L1o3S3RrQmV0Tjh5SXF0OHZmMXhRSUhNV0xVUzVZKytUTEhuampyNWRuNXpseApkdGcxT0lNNUI4em51azhXVmpGWFp0bkR0a0JsL1NDdXNML2R3MGtXY1lyVytJZkoyS3lIRTJTZlRKMmZIM01aCkpTdUV4VkNaRVBHZVlqSFN3TXd3elNCN0VUYWNxSVg4ZWVESmhVTmpPbVo3a0ZoTFByR0gzanpFSUlOaWpwV0wKQWdNQkFBR2pJekFoTUE0R0ExVWREd0VCL3dRRUF3SUNwREFQQmdOVkhSTUJBZjhFQlRBREFRSC9NQTBHQ1NxRwpTSWIzRFFFQkN3VUFBNElCQVFCc2toWE95ZDhXSjFOaVRDdWEvT0o4YXR6S3FiTzBoN205M2tKNWlWc2RZajRJCitKRFk3SHVPd1BqdC9sYXZrRVA2N1dsUWNPUlMrelpkVkRFUnVhVXhFU1FoREZWckNxeUV6MTY4WUErSG5WakwKcVZjb2kweURaVFV3a1BoTFJMNUYwTThzOWhuYTI2d2hLSng0WGlia0ZDMUx5VmFub0xOK1ZHOWhrdTZVZW01agpvbDRCVlN4bWM1b0V3QjNoN1J6ZE5OcE5vb3l1VUxlOTFDdzdISW52bU5PU1FMVkloeDdOeCtVWDY4aUx6NkR3CkEyS0Mwazczd1VFY2g3RXZOOW1aZGNKWEVSWnh0OTRGL3hVQStESmJkcWh1WGI0UzZnMS92cXdjQjY1Mm1oYTAKditreUR2RkVkQXJNeVdOeGthSHpYK1ZjQjFOT3hxY0xoLzNHZVUraAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    key: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBbTJnejc2cmVuWkxBZnRFcG9UNnBEamdiaWNVREluVWNTazJGNzF3cHFVaHR5bExDCkRPKytISUlpeThUR3pKL2d1Y1hZS21rOHdxR1R3ZW15dXVwT1NNTXA2TVFrcnd4aFhwWUxCVHZ1NHdEREJxbW0KQkszVTUvMndzSE9STmtlb21WMDZKOVhSNVpjK2RRUGZ5VkFITHptcGQwY2YyV0sxRzRVV1hmT1AyZXlyWkFYcgpUZk1pS3JmTDM5Y1VDQnpGaTFFdVdQdmt5eDU0NDYrWForYzVjWGJZTlRpRE9RZk01N3BQRmxZeFYyYlp3N1pBClpmMGdyckMvM2NOSkZuR0sxdmlIeWRpc2h4TmtuMHlkbng5ekdTVXJoTVZRbVJEeG5tSXgwc0RNTU0wZ2V4RTIKbktpRi9Ibmd5WVZEWXpwbWU1QllTejZ4aDk0OHhDQ0RZbzZWaXdJREFRQUJBb0lCQUZPUDBUVlNsRWI0RE1kago0bFdRWXNpQmhOVXNlUVlESUlZWGJ5Z0lUMkoxem9zV3VyN0gvbHBINHp2Yi8rVXhKbDNkc3VFREd5ZXdSOG5oCnhqZlpHdVRuQTliaitMR2pINHdEYzhPSnVXYVlGMFd5M05EeFEyVEd0VVg3cmg0WW8rQnJENFV4NUozbUdEQkYKT0FTQUlvelRIWHRFWkN5NGRaZHBsV0JKUVpVMFpVamhJK29SQWNpTHFyQ3BaaEVGbFVPRTFTTVM1S2lQc0NkRwp0Sjl2am1qN2cwZXJGM2RMNUxtVEhyc1lQSGM4cWpNZGpxdkd4c0RiaDIwSzQ1MlVTUDlxSDFNVDRhMEx3TDJnCkRxM2orb3dqY1VvRlAwdHBzbzljeElYRWl4UTZQSXpURHBvbzZOcC9WL3RrcnJlTmhMK2tXM0lDdXBPRTZldWEKcXpWUGxBRUNnWUVBd09LaFVRSDlubFBlRGJQRk00TmFIRVdjd3lTWldzeHFkV096cVhuaGErTEs1UlZYU3NNagpZcDlXZlNkNHkyYXB3dFVkZktDVnB5Mm9yc2d2bVoxUDJjOEtBbVd3K0s3eUQwOUNDVkFYemszdmhNcko0Z2ptCjI1RGJYREYvK1FsL0g2VGJmYWRQTWIwempEK0FyQ1oxU3RUVlBwTVF4RDlxdG1qTWZ3Y0FiZ0VDZ1lFQXprSWsKd3lGRnViUHREUGhBRDREa0QrZ3IzQjdGdEM5VWxrRUR6L1JtVmsyL2NGM0R6ZzFNME5IL28wNktkOE8wTlg4QgplcmNQZXlKTVRVSy9FcFR2S3pRSTR1L2g1ZVFpdHhBOHRCSDh2SHVwVXRvNENlOXU5bmN3Wm5LR3ZiOEdkSDlNClVBTGk5ekRoeUpuKzBtSjZyQU9KVHNjcVNIUHkxNmJwc3ljNDI0c0NnWUVBcGJZT3dYbmtXbEhUUkJKZUtaTi8KcHlwbk00QU9BR1ArVWp3RjdtUTN0bWh4eDc0OThJMFZxWFVhNlFzd1RBODNhWnVPYWJQTTNvUHJsNzJFcDRUdgpSVUVLYUdUVlZkRjNSSS9qTy8wRGRzcWVMSWZNU2RVOEFRYkNic0pZSDZ4NCtzYTNpNHhpRFdsdkQ5Nit2U2VOClBXejhoM1h4d2FoNkZaeVRrODZBSUFFQ2dZRUFzZHJYV0d1WWRFbHlYM3l4d0t5ai9CTjN2cGZLWTFWczJ5TzAKNWQwWllkSXBBZnZZbkJWYjU3VXRldVIvQWtiL1hpSG1aS3IxN25mazA3cDZpTXRrY1J5dGpRTE5DQzl3ditxQQpiY1lVNlhLNHozamNXYlRkT2lvTTBrcHZaYThUSWVHakxGdFEyMWFMV3k5dlRIc2V6TFUvOFc1TVI2MnorY0UwClJBZk9QNEVDZ1lBVU1BSEhlWjRpUG1lOGF3Z1ViWnRJWUZtV3BTQ1lPVytnQzVweGxQTmNaS3pDSmVrSWMvM0gKZzVlWEJLWE1XNXNTMWVWU1c2R0RKZWJxSkZkR3c2bmZBYi9RMzZ3TiswL0htMEdoNHJXY1JHSFJPUzdQaktzZAptdGVFbmUySjJyeEwrbFd5ZDRRYm82QVRyVHUzeGdEZUVxQW5ObXUrUWdkeU4zd1BrMTlHSVE9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
  version: v3.3.8
iaas:
  disk:
    size: 40G
    type: cloudssd
  image: abclid.vxd
  kernel: {}
  secret:
    value: {}
kubernetes:
  frontProxyCA:
    cert: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURHakNDQWdLZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREErTVNjd0ZBWURWUVFLRXcxaGJHbGkKWVdKaElHTnNiM1ZrTUE4R0ExVUVDaE1JYUdGdVozcG9iM1V4RXpBUkJnTlZCQU1UQ210MVltVnlibVYwWlhNdwpIaGNOTVRrd09ERXhNRGt6TmpFNFdoY05Namt3T0RBNE1Ea3pOakU0V2pBK01TY3dGQVlEVlFRS0V3MWhiR2xpCllXSmhJR05zYjNWa01BOEdBMVVFQ2hNSWFHRnVaM3BvYjNVeEV6QVJCZ05WQkFNVENtdDFZbVZ5Ym1WMFpYTXcKZ2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRRFJDK2ppb0FUZlVYVVFaUDU5eU5SUApDOTJud293VmN0Z0RBOU4rZGFLYkpsQVFKaitPdllEVlA0cVJaSFZQNDNXc2lZMXd6QUQ4MFhDMDhEcVRnU1NiCnFtZ1BtUE0xdjVXYlZJUGVrakRkT0pRbUZxU0pBRU96YnFwY2RjNlFhVWx3UEdXRi9XY3pjeHNkTjNxVWZZaXgKQWRvWFN2bEpnT1pVK01IeDNaNmNrREZoWUlpbXdqUE5wUFJOMEYxQm1TRDh6TzhFeWUvTkJDY1orTWJSV0QxRQpobGpMYVNucWo0V1FZdzl3a0NjQ2dhWUF3K1dRUEZSVlZyUE1samNXelJYTFZLbjlhcnV3Umc2NzEzajhMRHdpCnB6eHlaWTZnYitoVjNzZzNBMkQrR3VyRXJvK0laN3FLUGxuTmhFRjZwS01JUnNCU3ZFMldEYlpiL2hHY3N6WkYKQWdNQkFBR2pJekFoTUE0R0ExVWREd0VCL3dRRUF3SUNwREFQQmdOVkhSTUJBZjhFQlRBREFRSC9NQTBHQ1NxRwpTSWIzRFFFQkN3VUFBNElCQVFCUWxDM05odXZzdkZjZWtoeHJnSTFUUExDcStkTzlVVjRkWjBYR3N6SXM2SFZRCjJqTGsxVkQ2NTJnVGJ4U05hdDZSMG01ZnhMVUlNYmVUd2FHcUR4U3lKOCtueDFGMWN0aG9TYkoyZVRuOXhhYWsKR2xwRno3ZXJYLzRjQTNLVVd3ZUFKVkd5Q1JaOUs3bkswNUVpdk9CWC95SmZGUGl1MXROOTVIT0tsZmNUTUNuVwpqdGppRDc3SlNmdDRMcW9VM0h1cUhmVUdENElnUTJNL3Z0bndBTjFKdWhaZmJPcnNBM1J6T1pkcWJTMXFnMXNxCkJFYjErdCtuMkQyNVZxa2NXeWJIRlZNTlhUajVuNmhjaG1Scmw0aG83RG4xQVFxSEVSN05xVHBVMHpmYWhCQ3cKY3ZaV3VWMnhxOWppNFBzZkFnekZldmFyT3NVSHB2WG1tdWhEbWF1UQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    key: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBMFF2bzRxQUUzMUYxRUdUK2ZjalVUd3ZkcDhLTUZYTFlBd1BUZm5XaW15WlFFQ1kvCmpyMkExVCtLa1dSMVQrTjFySW1OY013QS9ORnd0UEE2azRFa202cG9ENWp6TmIrVm0xU0QzcEl3M1RpVUpoYWsKaVFCRHMyNnFYSFhPa0dsSmNEeGxoZjFuTTNNYkhUZDZsSDJJc1FIYUYwcjVTWURtVlBqQjhkMmVuSkF4WVdDSQpwc0l6emFUMFRkQmRRWmtnL016dkJNbnZ6UVFuR2ZqRzBWZzlSSVpZeTJrcDZvK0ZrR01QY0pBbkFvR21BTVBsCmtEeFVWVmF6ekpZM0ZzMFZ5MVNwL1dxN3NFWU91OWQ0L0N3OElxYzhjbVdPb0cvb1ZkN0lOd05nL2hycXhLNlAKaUdlNmlqNVp6WVJCZXFTakNFYkFVcnhObGcyMlcvNFJuTE0yUlFJREFRQUJBb0lCQUZ2RFBwYzhadWNnZXFLWApzcDdFYVczSlo2TWNZeUdIS0FzcXdzdmdGMkREa0tHR0tLQWZ6MDNNZHFjYjBlTWZsYWdIT1c4cUhjVGNxdnpCCjl2U0kyK2o1QkhUVVR1NXBDdU1FVmw1OURiWU4vL280TmtGdFBFcW5hV0RzMVovT2w4NE10UVA4R0RFZGRlbDkKVVBHZHFVTUo3UklNZHlFczArNjR3Mm5JUHJlNlNOb003eFJMUUxzcHhMV3plUUptMU9XQS9WNHhrZUgrb29rMQorOThSUXM0VEZzb0hCREhzanQzLy92cnhzWmJqeUJkdEN3cTQyeElmSWh5dmNnMkUzRS9SZFI4ZEVWb1JWSTRUCmdtSnVKRlJLNStMaEVwNUd2S3hxYXhCa0hmekJGTmlzbTdYcWwyVWF0MG5YbUZ6OFpQMDNTZzAvMGNGZkREc0oKSXVFWnZJRUNnWUVBNXFXNXA0UkJCbWxwSjE5WklzWnZkT1Z0Qmd2S0drMjA1OGFZb2JycDZjKzkyMFdiazVzaQozVzQ1ekhaNjQ5YkdxWXYzMEM4RHJIRzBoQkxYVWZieHY2YmVGMXUvREtGRDl6NFY5OTdPUXQ2MVJ0WWprTUFsCm43Mi9rWEFuN0h0eDNBdDlYalZoWktmRVR2TVU4eGJENGF0cHk5N1UvNXFxMVFRV2l0M3lMZ3NDZ1lFQTZBWloKb01TVmZqVGVyMUl6bS9XY2doOG5ONUxzQ3R1MXVKQmhoZ0QxQUUvYXEvYmZab2lESlZXU0ZKc2IvREJobEkrcApVaHRyazVMUUl6dTVvMThua0o2QXg3djkxVXVqYThrbFR4WVoydzhnK3AwSVIxbW9UTEg0Mk5vVi9PY3djbkg5CkwwNzdaV2duTlliZUR2MTJ6b1REUXdHZ2FIUkNsTXVjcHgzazd1OENnWUJ5Zk9McC95RVQ0TEVjcFJ5bXdWOUQKNURvNDNSTkYrVHFLTGk5SFlIT0o1dCt6L3hwWnE1RWozdm43dnZnRExuSlFhTFRxOXR0WTl0d0hEeDhvaFc2OApsa1Q1elVYSkxDZURpNkwxOWZmbWc0dnlESXQ4NTVRRmRmZW9ac2E2Z1JBa2pPTi9Kdm9nTDVLbktjeEZXaENECmJVWEh0K1Y4dHphREpGTllQUkZndndLQmdRQ1JnV0gxd3pKbzJpa0lVNG1QOWFBM3JlZS9IMEV0c2drcy9FWmMKYmY4M09kek5XTjFTaEt1UjN5N2tBejJ5a25pdlhNUjNmRUNqWkQ5b3lReXEyb0tLWEF3d3RjRUNZUlBVQldtRApSajNpdFlNZUJ4cG8vRjNoOHY0MnA5V0FLMCtqaGI2Y1MzQzJjSEdlVEx5M005YXN2bTloZHdTc1hMUmdjYXdFCjFXZDdPd0tCZ1FDUjBJY09EWW9FSWVHenpMWFJMczkrd2FWNG9hY1ZTandNOStvNVZVNEJBeklKWTM1RDd4c0IKZGdwbGRCVHpsVG16WEo4NUJQTU0yN0E5SGtrTlN5QS9SekQ1MFdKcnlTQkViODNmdzNaVG9VMVlWZWZiOWtjdQpHNFFpQ3FpOU9Oby9hcm43ZkFzS29tcUluVFFnQ2VmSFpxSVJuMHdGV1ZqbHpDUFVla3BwdVE9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
  kubeadmToken: kbcg7w.bmyrys20y0clkp3v
  name: kubernetes
  rootCA:
    cert: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURHakNDQWdLZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREErTVNjd0ZBWURWUVFLRXcxaGJHbGkKWVdKaElHTnNiM1ZrTUE4R0ExVUVDaE1JYUdGdVozcG9iM1V4RXpBUkJnTlZCQU1UQ210MVltVnlibVYwWlhNdwpIaGNOTVRrd09ERXhNRGt6TmpFNFdoY05Namt3T0RBNE1Ea3pOakU0V2pBK01TY3dGQVlEVlFRS0V3MWhiR2xpCllXSmhJR05zYjNWa01BOEdBMVVFQ2hNSWFHRnVaM3BvYjNVeEV6QVJCZ05WQkFNVENtdDFZbVZ5Ym1WMFpYTXcKZ2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRRER0UjRoZDAyczNoK1lOTzZ5MzFJMQpxS2orWlJxL0s2UWVwKzBaMmF5KzAwMnhKYXdxQmM4a3ovWjlkbEN3Z1JRYVFZQnBHblIwMndVTWxZb1RTcG1JCmVwRWg3elJjU2I1UXo5L1VvRUlFTDd3NEFtVVRnSFdhTWs3SHFjSVBjaUhKLzA0MVZyRkQrVlo2NmxkWHVOOHQKYkxWWENQTXE1QWhuTXpBT05CcXBGNDA5UlJaVSt3T3JJdVdxVDlnOXVRcWJXOHBSeHJJaVNsUm9PODl4RnluMQp6TTZDcittMUt4YllDbm14U2c2K2FRaE9EQnZGbEVUYTZJUTNkMElCakkvVHRHNXNTU0RjSDNRcWcrRDZKUTFEClhYbDlsSC9ZUE40b1pxTWRRVFZKSXVNVU5EZm5zUDJoVGlJZ1Y4c20xTDJkVm1NR1Q5RHl2b0ZFVGl6YU5XL1QKQWdNQkFBR2pJekFoTUE0R0ExVWREd0VCL3dRRUF3SUNwREFQQmdOVkhSTUJBZjhFQlRBREFRSC9NQTBHQ1NxRwpTSWIzRFFFQkN3VUFBNElCQVFBc2IyMGdOWndIV0JLV0JEYWtrdHk2bXgyaFhKcVk2Wm9TUUJueTNFY2k4U1V3CjVORHoveFROZ2FtdXFVajlGVkNRQ0w3c1psQkM2ck9SbFlQcDJWT3BJOXZUbzRRa3I5Q1doY3grK0V1dXZWZmUKMVdwYnc2RUV0eXZMSkNXdklVVDZrOTdJWVV5SEhsM013ZUE2SytvUG5FV2xibGRjdllqR1M0TWI5WXBKVk1zKwpQSURBaEpVYzZscnVvN2FUNEx3a1lUNFNEZm0rQ3M3UUtNOUNiakRXdTNQbFJmdlBrSHpLTE5XYUc4aWVPOGNKCmNpb3crSUdlUFo2SGN1cE9DZUE3c0E0VTV4TG56UXFLeHMxQlI5TXVSZEppcU45WDN4aVdtSWFsYzJUUDJpaWEKVWt6UmJ3SHlUSFhsMi9XT2ZrWkw3NjB5WlV0dFNFT0xZY29ncklxbgotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    key: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb2dJQkFBS0NBUUVBdzdVZUlYZE5yTjRmbURUdXN0OVNOYWlvL21VYXZ5dWtIcWZ0R2Rtc3Z0Tk5zU1dzCktnWFBKTS8yZlhaUXNJRVVHa0dBYVJwMGROc0ZESldLRTBxWmlIcVJJZTgwWEVtK1VNL2YxS0JDQkMrOE9BSmwKRTRCMW1qSk94Nm5DRDNJaHlmOU9OVmF4US9sV2V1cFhWN2pmTFd5MVZ3anpLdVFJWnpNd0RqUWFxUmVOUFVVVwpWUHNEcXlMbHFrL1lQYmtLbTF2S1VjYXlJa3BVYUR2UGNSY3A5Y3pPZ3EvcHRTc1cyQXA1c1VvT3Zta0lUZ3diCnhaUkUydWlFTjNkQ0FZeVAwN1J1YkVrZzNCOTBLb1BnK2lVTlExMTVmWlIvMkR6ZUtHYWpIVUUxU1NMakZEUTMKNTdEOW9VNGlJRmZMSnRTOW5WWmpCay9ROHI2QlJFNHMyalZ2MHdJREFRQUJBb0lCQUdyNzdaYTUwenAyeXFxMwo5T3pQYSs5dFhjU0RuSlY0MCtUMlE0VG9HNmpOZGlFcXlPekg5QzdaaUpPWlJBN1k0UlpoVEVNb2RSVVVUYlJOCll3ZWN0a1dIU3lOVDBqbkpEa0s3QUU4SnRFejVrMWpDNW1JOHpRMFlCenphcmYwbmxSVUpmY1ZtdU16QjF1YW0KaUV5cUFVYWhzSmY5aW9DZDI0SWUxMXhNVzI2blBnS2cwdzMyMXlyR1MwaDBYZVAybmFmNTEwbFAzTFljWkhCdQpESEJCbWZLUzRNRjduTnFOS2oxRW8rSGxyeERuL2ZPN2YrTjRReTlXd3FZZERlWHQ4aU44NnZ6Q09CRC9wZnAyCk4vVFM3ekZ1UzVmMjFQaDJJY2pMZU9MYjEwUU9rVnlDZk9FRGtEbTRycy9wb2NCUUN5VlBORnd0bWtvTWp5Z0kKR1lLQ2NURUNnWUVBMFJwT2VqdDZrNDArdUhaSDFqL2FPalUyMlBWRGpLMXExZDNHNkJrVE01YytwL3NaVjZKYQpDeEY1MVU5MmpwS1lNdmdRMWZhU0JDMkg3TlpxMDRFMjVGVVE1clZuczg3cE05WEg0VmcyL0RQb2ZhMjcwUXJPCmFLeXlMYVRnUkhDTmE0NldIdVY2dmUxZ1dsMUhMb0NrUGJpeUh2SGQ0NWtETm5JbnBONk8xc1VDZ1lFQTc1bTMKME9va0gxU3dJSlFraTRkZEFra2k5aXRYanNUd2RFSndkNWtWWGlNZ1RGQXhSWC9qTzdEMjdsWDBCVkdteStmVQpBTDhFQ00wWEVBeHVrZVRSWC9nZkhpenJCaHBLWGRBcXVYdkFaRHV0aXBEWTYwbk0vbnk5UEdyQS9lRmpjNVdLCkJkVlIzWFNPQXJIdExBRnRNdGJPSUtSUTF1YkFoMC9PRHBwSTFiY0NnWUIzdFZhK05YVHNLZzJCTHYzYlV6ek0KM0JBbFR1dzRDa1BDWWkvd1NnS3JJMmdVWlBWU0xUamRZMGpiYkoyVDY2ckVheTNBUUdQQmpvdGxkQzgrSEpoZgpTYS9lVkhZbEEweVFoMC9oMjAzSFByUlgxdkZTTUp5UVltV3pLZFBXZXBVTHdWcUNINkFRUVdoSzgyRy81cGVnCldpOW05ZEt3N0xWaEl5TTlDTWkvZlFLQmdDRXZaRE11U3NTQUlValIyK0hyWktsdFljZEFwNGJocGdBa041bUkKL1ZtbGVkSzRCS3NBbElOdTlqUjZlU3JIYlZldWkwRnpNMmJZcVFvUy9ybDhQTGVURVJSSDJtRmxBTW5QakN0VwpoWVljY2VQUVBnc1FsTERtcS9zVE9obXZ2dXVDV2JTSElDaVEvUjVJY0hYNDJKd0MvbGV0Q25sSlNuOHpRWEhxClVvNHZBb0dBRXdrSnI5K2FkYnFudnF5aWhoRFpwV1NFK1djWFJyakpMV2lsVzJ4c2hZbC9TOUZwUG1senlYMWUKOUovbHUzRWhhcDFPRFl1SVpMM0JINHN4c3E5dHNseWRwTXBWeFhkQVdhc1AvWkRBUSttZFQ0S1RNSWRUUWxGcwoyWXN6dDdiQkdadHowbjZSd0JmVU4zeHpBaGF0N2R5a0N2bDFodUpaZWZUeWlUTTZBYzA9Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==
  serviceAccountCA:
    cert: LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUFvaEwwMUU3ZVI2SWR5bnRVVHlucQpFeVRvay9WWEhMS0ZDQ2piT1JxMndLRGE4SU5yQng4OUNZUW9TdnBwN2JMMkdTdmZZU2tReDArVUMxUWF4Q0IxCkJEQmZzQndMVzdwcmx4QUd0UnNabTVINE1abUhBY0hZR0k4anJzeFdoOWhBRU40bWlTUVRZK1NMamYybStTZEcKdUFRaGtIeG1Vd0Z6WWVxTEU3UDhyNnZ3REpML3Uva2pqdytPY3Z5WldGd2tCTyt0TEY3T1pMMmpMQ0FnMi94cwpnMitDR3cvM0N4dUt1L1BVU0dYeWtFQitZSERrcDZoK2V0RHYwTDRCTW1pTGlGaTNFQndvM0hHL3dHZG83ZlNjCnFzSG9Nc25mNjdBTHQ1eTMrNnU4UEE2Y2lsQi92WDR1bVZLUlc4UUxuY1VjZ3pyTFNjSjVpdjA5MWVHQm96dUsKa1FJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg==
    key: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcFFJQkFBS0NBUUVBb2hMMDFFN2VSNklkeW50VVR5bnFFeVRvay9WWEhMS0ZDQ2piT1JxMndLRGE4SU5yCkJ4ODlDWVFvU3ZwcDdiTDJHU3ZmWVNrUXgwK1VDMVFheENCMUJEQmZzQndMVzdwcmx4QUd0UnNabTVINE1abUgKQWNIWUdJOGpyc3hXaDloQUVONG1pU1FUWStTTGpmMm0rU2RHdUFRaGtIeG1Vd0Z6WWVxTEU3UDhyNnZ3REpMLwp1L2tqancrT2N2eVpXRndrQk8rdExGN09aTDJqTENBZzIveHNnMitDR3cvM0N4dUt1L1BVU0dYeWtFQitZSERrCnA2aCtldER2MEw0Qk1taUxpRmkzRUJ3bzNIRy93R2RvN2ZTY3FzSG9Nc25mNjdBTHQ1eTMrNnU4UEE2Y2lsQi8Kdlg0dW1WS1JXOFFMbmNVY2d6ckxTY0o1aXYwOTFlR0JvenVLa1FJREFRQUJBb0lCQUMwMWpuZStmRUNORnpSSgpESzU4YVovbmRyejIvZEt5ZFd0ZVFqQitwQ1c5ZVBVSUk0MjhDQUdrakx6RmQxRG56OEFidmRiVFpxdkNKZkMzCkNEUlhQV3pxdFhaaGpFV05EYi80cDNaVFZlUlFabVFuaUVKVU9SdzRxV0p6KzFzdlZrZWVRQytYRVpXV3hkclkKdTMvYWxkNTB3SXdXTjFER0lkSlpZMndlazhqSWdxSlVQUk92bi9CYUl3TWcyRHRnMTlYb0lDZWsyVXpkVWpPTwo1T3Q1RUNGRXdPeGZxUXI1VEhFYWNPekdNajlDa0RPdjl6Qis3cUtqRXRRYjRIS01HMkJRYzR4OHhpNVpVYm9NCklHTmtzVVE4a0FlU21haEcwOXcvNjhrdjlRRkt0VEdjbVpsZXNWbDhWV1h0UHEwS0Q4c2hOSVNWbjU2RFZ4Ni8KZ1lhWE9GVUNnWUVBeGtiYjFpZTZTZ3JqRmJuR2NadUpLSkU1QVVrZHJ0ZCtGRU04QWM5Z1pia3FzODB6MklxVwp5MlFKRGRJUzRLTVh2eDdScFpqK1paRmtVQ1hmcU9kMElkUUtRRW5CTVEzMGRTRFhBR3J6NzY4REZockR2Q2w4CklvSXgwR3l2SVBmd2FCWlNDcFdYN2JaaG0wMlcvOUFJTWc3d3RmaXdEY1RCOE4rREpZcndnWk1DZ1lFQTBVSDgKUEk0V2pYTHl2bHdXN0JDcHA1THFCWldGTWhHTHNlZ01ydHgzSGxxQjhXVTVRK3kvWnhuZVpVbkRhQnk3YWhqTQpsM3ZPajVXZDVWak1LUkVSN3Yrd3VxM2ZYTjV5VDhXQWZubmltY2tSbDZCNjkvVUFSeG9Pdnd0cGxIVHozWHg0CkUzZVRiUkM4ektFSmJTNUlVdmw2VFhsZlUzSzZpdVlyOFdMUzZjc0NnWUVBcnNFUlRVVWltZFBTcGZsaEFBeVgKN050aUpOSHVpWVdBcUJkQ01rNDJwakYzZzVXZTFvSC91aS9uRXZsT2poTHhBUkFHc3krUE9MSFdlaFdIZFhUYQpGRjZ3MGt5dks2OGpBSUQ3UG5FRm93RGJkWVlOY2pBV0tzd1pYNXdMRnNHd2IrME1UaXZmQmpLekFKQjRQK3Q0CjdiWGhUbUZydWhicTRJUC9NUDJ1VUdVQ2dZRUFsd3pGa2NTVEhQbXJwYU45M0Z0T3k2cXVDT2ZjWkw3cmtybEoKWm5PMy9JNGluR2lRQktzQm90KzJmSERaZis2MWppbG1qYmFON1hGM3I1VUFrbWhEQkwxSENnbjJZT2dscGRXUQpJanZEU1hVdG9NRHo0c2JVczM1b3hKanRWbjl4aFNDUzRLS0JKY3BlTG12VURSN0trREtMaVI4aW8yNytuc0wxClYySVpreThDZ1lFQXJpVlNOOFFxektsalhPWWNEZ1dHbW1LM3FsajUvRVRGckppMllFTHVBSGw1TEJzR1RHVDMKbFNwS0Z1dHh1bEFoWXVrSElXUHNsMWxqSjlCOURjTnJ5RmlVM3JtRDhpVFhVWjIzMTAwcjM5NWVFLzBDY3MxTApTSUNNU2M2bVFtQks2Szl2WWtwdFRhL2NmNU8zZ24yeG9IODYydnJYVHJDTzhobDQ0eDg1Y3ZvPQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
  version: 1.12.6-aliyun.1
namespace: default
network:
  domain: cluster.domain
  mode: ipvs
  netMask: "25"
  podcidr: 192.168.0.1/16
  svccidr: 172.10.10.2/20
registry: registry-vpc.cn-hangzhou.aliyuncs.com/acs
sans:
- 192.168.0.1
`

func TestCCMAuthConfig(t *testing.T) {
	cfg := &v12.Cluster{}
	err := yaml.Unmarshal([]byte(bootcfg), cfg)
	if err != nil {
		panic(fmt.Sprintf("%s", err.Error()))
	}
	_ = &v12.Master{
		TypeMeta: v1.TypeMeta{
			Kind:       "NodeObject",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name: "mynode",
		},
		Spec: v12.MasterSpec{
			ID:   "",
			IP:   "",
			Role: "Master",
		},
		Status: v12.MasterStatus{
			BootCFG: cfg,
		},
	}
	actx := actions.NewActionContext(nil)
	err = NewActionCCMAuth().Execute(actx)
	if err != nil {
		panic(fmt.Sprintf("%s", err.Error()))
	}
}

func TestAnonyStruct(t *testing.T) {
	tpl, err := template.New("wdrip-file").Parse(wdripf)
	if err != nil {
		panic("failed to parse config template")
	}
	// execute the template
	var buff bytes.Buffer
	err = tpl.Execute(
		&buff,
		struct {
			Version  string
			Registry string
		}{
			Version:  "0.1.0-149bf7ce",
			Registry: "registry-vpc.cn-hangzhou.aliyuncs.com/acs",
		},
	)
	if err != nil {
		panic("error executing config template")
	}
	fmt.Printf("%s", buff.Bytes())
}

var wdripf = `
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: wdrip
  name: wdrip
  namespace: default
spec:
  ports:
    - name: tcp
      nodePort: 32443
      port: 9443
      protocol: TCP
      targetPort: 443
  selector:
    run: wdrip
  sessionAffinity: None
  type: NodePort
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: wdrip
  name: wdrip
  namespace: kube-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: wdrip
  template:
    metadata:
      labels:
        app: wdrip
    spec:
      priorityClassName: system-node-critical
      containers:
        - image: {{ .Registry }}/wdrip:{{ .Version }}
          imagePullPolicy: Always
          name: wdrip-net
          command:
            - wdrip
            - operater
            - --bootcfg /etc/wdrip/boot.cfg
          volumeMounts:
            - name: bootcfg
              mountPath: /etc/wdrip/boot.cfg
      nodeSelector:
        node-role.kubernetes.io/master: ""
      tolerations:
        - operator: Exists
      volumes:
        - name: bootcfg
          configMap:
            # Provide the name of the ConfigMap containing the files you want
            # to add to the container
            name: bootcfg

`
