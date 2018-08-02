class OutlyerCli < Formula
  desc "Outlyer CLI allows to easily manage your Outlyer account via command line."
  homepage "https://www.outlyer.com/"
  url "https://github.com/outlyerapp/outlyer-cli/releases/download/0.2.0/outlyer-cli_0.2.0_Darwin_x86_64.tar.gz"
  version "0.2.0"
  sha256 "4bbb0307d3a96144fa3e0572d98de3aafb382c6bb407be20cee91d3652a6a803"

  def install
    bin.install "outlyer"
  end
end
